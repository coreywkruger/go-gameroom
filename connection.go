package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Connection -
type Connection struct {
	ID      string
	send    *chan []byte
	receive *chan Message
	kill    *chan bool
}

// Message -
type Message struct {
	ID      string
	message []byte
}

// ConnConfig -
type ConnConfig struct {
	send    *chan []byte
	receive *chan Message
}

// NewConnection -
func NewConnection(id string, cfg ConnConfig) *Connection {
	send := cfg.send
	if send == nil {
		ch := make(chan []byte, 256)
		send = &ch
	}
	receive := cfg.receive
	if receive == nil {
		ch := make(chan Message)
		receive = &ch
	}
	kill := make(chan bool)
	return &Connection{id, send, receive, &kill}
}

// Start -
func (c *Connection) Start(ws *websocket.Conn) {
	// read/write loops to websocket; read bocks until closed
	go c.Writer(ws)
	go c.Reader(ws)

	<-*c.kill
}

// Reader - reads message from websocket; puts in `receive` channel
func (c *Connection) Reader(ws *websocket.Conn) {
	for {
		_, Msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Couldn't Read:", err.Error())
			*c.kill <- true
			break
		}
		*c.receive <- Message{c.ID, Msg}
	}
	ws.Close()
}

// Writer - writes message to websocket
func (c *Connection) Writer(ws *websocket.Conn) {
	for Msg := range *c.send {
		err := ws.WriteMessage(websocket.TextMessage, Msg)
		if err != nil {
			log.Println("Couldn't Write:", err.Error())
			break
		}
	}
	ws.Close()
}
