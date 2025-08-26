package webrtc

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Hub  *Hub
	Send chan *Message
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		//var msg Message
		var msg map[string]interface{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket.IsUnexpectedCloseError: %v", err)
			} else {
				log.Printf("error: %v", err)
			}
			break
		}

		//log.Println("Received msg: ", msg)
		//msg.SenderID = c.ID
		//c.Hub.Broadcast <- &msg

		log.Println("Received msg event: ", msg["event"])

		message := &Message{
			Event:    msg["event"].(string),
			Data:     msg["data"],
			TargetID: msg["targetId"].(string),
			SenderID: c.ID,
		}

		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
