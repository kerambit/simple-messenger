package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kerambit/simple-messenger/internal/webrtc"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// add domain validation in prod mode
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *webrtc.Hub, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("userId is required")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &webrtc.Client{
		ID:   userID,
		Conn: conn,
		Hub:  hub,
		Send: make(chan *webrtc.Message, 256),
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
	hub := webrtc.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	port := "9000"
	log.Printf("🚀 Server starting on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
