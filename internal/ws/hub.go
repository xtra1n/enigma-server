package ws

import (
	"encoding/json"
	"enigma-server/internal/enigma"
	"strings"
)

type Hub struct {
	enigma *enigma.Enigma
}

type Message struct {
	Result   string `json:"result"`
	Position string `json:"position"`
	Error    string `json:"error,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		enigma: enigma.NewEnigma(),
	}
}

func (h *Hub) Run() {
}

func (h *Hub) HandleMessage(msg string) string {
	msg = strings.TrimSpace(msg)

	parts := strings.Fields(msg)
	if len(parts) == 0 {
		return h.errorResponse("empty message")
	}

	command := parts[0]

	switch command {
	case "set_positions":
		if len(parts) != 2 {
			return h.errorResponse("usage: set_positions AAA")
		}
		if err := h.enigma.SetPositions(parts[1]); err != nil {
			return h.errorResponse(err.Error())
		}
		return h.successResponse("OK")

	case "set_plugboard":
		if len(parts) < 2 {
			return h.errorResponse("usage: set_plugboard AB CD EF")
		}
		plugboard := strings.Join(parts[1:], " ")
		if err := h.enigma.SetPlugboard(plugboard); err != nil {
			return h.errorResponse(err.Error())
		}
		return h.successResponse("OK")

	case "get_position":
		pos := h.enigma.GetPosition()
		return h.messageResponse("", pos)

	default:
		result := h.enigma.TranslateString(msg)
		return h.messageResponse(result, h.enigma.GetPosition())
	}
}

func (h *Hub) successResponse(msg string) string {
	response := Message{
		Result:   msg,
		Position: h.enigma.GetPosition(),
	}
	data, err := json.Marshal(response)
	if err != nil {
		return `{"error":"JSON marshal error"}`
	}
	return string(data)
}

func (h *Hub) messageResponse(result, position string) string {
	response := Message{
		Result:   result,
		Position: position,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return `{"error":"JSON marshal error"}`
	}
	return string(data)
}

func (h *Hub) errorResponse(err string) string {
	response := Message{
		Error: err,
	}
	data, _ := json.Marshal(response)
	return string(data)
}
