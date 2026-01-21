package main

import (
	"enigma-server/internal/ws"
	"flag"
	"log"
	"net/http"

	"github.com/xtra1n/enigma-server/internal/ws"
)

func main() {
	port := flag.String("port", 8081, "Server Port")
	flag.Parse()

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":" + *port
	log.Printf("🚀 Starting Enigma server on http://localhost%s", addr)
	log.Printf("📡 WebSocket endpoint: ws://localhost%s/ws", addr)
	log.Printf("❤️  Health check: http://localhost%s/health", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
