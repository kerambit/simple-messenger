package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kerambit/simple-messenger/internal/webrtc"
	"log"
	"net/http"
	"os"
	"path"
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

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cwd, _ := os.Getwd()

	htmlPath := path.Join(cwd, "static", "index.html")

	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		http.Error(w, "HTML file not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, htmlPath)
}

func main() {
	hub := webrtc.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/", serveHome)

	port := "9000"
	log.Printf("🚀 Server starting on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
