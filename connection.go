package main

import (
  "log"

  "github.com/gorilla/websocket"
)

// Connection contains websocket connnetion
type Connection struct {
  ws       *websocket.Conn
  outbound chan []byte
}

// CreateConnection -
func CreateConnection(w *websocket.Conn) *Connection {
  return &Connection{
    ws:       w,
    outbound: make(chan []byte, 256),
  }
}

// Init -
func (c *Connection) Init(broadcast chan []byte) {
  go c.Writer()
  c.Reader(broadcast)
}

// Reader -
func (c *Connection) Reader(input chan []byte) {
  for {
    _, Bytes, err := c.ws.ReadMessage()

    if err != nil {
      log.Println(err.Error())
      break
    }
    input <- Bytes
  }
  c.ws.Close()
}

// Writer -
func (c *Connection) Writer() {
  for message := range c.outbound {
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Println(err.Error())
      break
    }
  }
  c.ws.Close()
}
