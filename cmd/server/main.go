package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"enigma-server/internal/ws"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	hub := ws.NewHub()
	go hub.Run()

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	exeDir := filepath.Dir(exePath)
	clientPath := filepath.Join(exeDir, "..", "..", "client.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, clientPath)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleWebSocket(hub, w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":" + *port
	log.Printf("Starting Enigma server on http://localhost%s", addr)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("Client HTML: %s", clientPath)
	log.Fatal(http.ListenAndServe(addr, nil))
}
