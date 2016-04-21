package main

import (
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

var broadcast = make(chan []byte)

// Register blah
func Register(w http.ResponseWriter, r *http.Request) {

  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  m, registerErr := RegisterMember("123", &Connection{
    ws:       ws,
    outbound: make(chan []byte, 256),
  })

  if registerErr != nil {
    http.Error(w, registerErr.Error(), 500)
    return
  }

  m.WS.Init(broadcast)
}

func main() {

  // Main http rounter
  r := mux.NewRouter()
  r.Headers("Access-Control-Allow-Origin", "*")
  r.Headers("Access-Control-Allow-Headers", "Content-Type")
  r.Headers("Content-Type", "application/json")
  r.Headers("Content-Type", "text/plain")

  // GET Join room
  r.HandleFunc("/register", Register)

  go func() {
    for {
      select {
      case output := <-broadcast:
        log.Println("OUTPUT", string(output))
      }
    }
  }()

  err := http.ListenAndServe(":3334", r)

  if err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
