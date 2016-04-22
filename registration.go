package main

import (
  "math/rand"
  "net/http"

  "github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

// CreateRegistrationHandler - registers new member & connection
func CreateRegistrationHandler(receiving chan []byte, closing chan int) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    id := rand.Intn(10000)
    m, regErr := RegisterMember(id, &Connection{
      ID:      id,
      ws:      ws,
      send:    make(chan []byte, 256),
      receive: &receiving,
      closed:  &closing,
    })

    if regErr != nil {
      http.Error(w, regErr.Error(), 500)
      return
    }

    m.Connection.Listen()
  }
}
