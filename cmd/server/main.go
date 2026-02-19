package main

import (
	"flag"
	"log"
	"net/http"

	"enigma-server/internal/ws"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	hub := ws.NewHub()
	go hub.Run()

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

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
	log.Printf("Serving static files from: web/")
	log.Fatal(http.ListenAndServe(addr, nil))
}
