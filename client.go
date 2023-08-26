package main

import "github.com/gorilla/websocket"

type Client struct {
	socket  *websocket.Conn
	receive chan []byte
	room    *Room
    clientName string
}

//The client is reading a message from the user
func (c *Client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}
        messageWithName := c.clientName + string(msg) 
		c.room.forward <- []byte(messageWithName)
	}
}

func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			return
		}

	}
}
