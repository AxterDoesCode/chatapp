package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
	forward chan []byte
}

func newRoom() *Room {
	return &Room{
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		forward: make(chan []byte),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (room *Room) addClientToRoom(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r , nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{
        socket: ws,
        receive: make(chan []byte),
        room: room,
        clientName: fmt.Sprintf("User%d ", len(room.clients)),
    }
    room.join <- client 
    go client.read()
    go client.write()
}
