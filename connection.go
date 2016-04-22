package main

import (
  "github.com/gorilla/websocket"
  "log"
)

// Connection -
type Connection struct {
  ID      int
  ws      *websocket.Conn
  send    chan []byte
  receive *chan []byte
  closed  *chan int
}

// Listen -
func (c *Connection) Listen() {
  defer func() {
    *c.closed <- c.ID
  }()
  go c.Writer()
  c.Reader()
}

// Reader -
func (c *Connection) Reader() {
  for {
    _, Bytes, err := c.ws.ReadMessage()

    if err != nil {
      log.Println(err.Error())
      break
    }
    *c.receive <- Bytes
  }
  c.ws.Close()
}

// Writer -
func (c *Connection) Writer() {
  for message := range c.send {
    log.Println("received this: ", message)
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Println(err.Error())
      break
    }
  }
  c.ws.Close()
}

// Kill -
func (c *Connection) Kill() {
  close(c.send)
}
