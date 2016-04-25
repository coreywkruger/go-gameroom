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
}

// Message -
type Message struct {
	ID      string
	message []byte
}

// Reader - reads message from websocket; puts in `receive` channel
func (c *Connection) Reader(ws *websocket.Conn) {
	for {
		_, Msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Couldn't Read:", err.Error())
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
