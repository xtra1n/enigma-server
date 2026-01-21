package ws

import (
	"encoding/json"
	"enigma-server/api"
	"enigma-server/internal/enigma"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingInterval   = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	enigma *enigma.Enigma
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		enigma: enigma.NewEnigma(),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket Error: %v", err)
			}
			break
		}

		var msg api.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			c.SendError("Invalid JSON")
			continue
		}

		switch msg.Action {
		case "set_rotors":
			// Конвертируем слайс в массив [3]string
			if len(msg.Rotors) != 3 {
				c.SendError("Need exactly 3 rotors")
				continue
			}
			rotorArray := [3]string{msg.Rotors[0], msg.Rotors[1], msg.Rotors[2]}
			if err := c.enigma.SetRotors(rotorArray); err != nil {
				c.SendError(err.Error())
			} else {
				c.SendOK()
			}
		case "set_positions":
			if err := c.enigma.SetPositions(msg.Positions); err != nil {
				c.SendError(err.Error())
			} else {
				c.SendOK()
			}
		case "set_plugboard":
			if err := c.enigma.SetPlugboard(msg.Plugboard); err != nil {
				c.SendError(err.Error())
			} else {
				c.SendOK()
			}
		case "encrypt":
			result := c.enigma.TranslateString(msg.Text)
			c.sendResult(result)
		default:
			c.SendError("Unknown action: " + msg.Action)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendOK() {
	resp := api.Response{
		Status: "ok",
		Pos:    c.enigma.GetPosition(),
	}

	data, _ := json.Marshal(resp)
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) SendError(errMsg string) {
	resp := api.Response{
		Status: "error",
		Error:  errMsg,
	}
	data, _ := json.Marshal(resp)
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) sendResult(result string) {
	resp := api.Response{
		Status: "ok",
		Result: result,
		Pos:    c.enigma.GetPosition(),
	}
	data, _ := json.Marshal(resp)
	select {
	case c.send <- data:
	default:
	}
}
