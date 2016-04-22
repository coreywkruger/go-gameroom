package main

import (
  "log"
  "math/rand"
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

// IncomingMessages - channel used for incoming messages from ws
var IncomingMessages = make(chan []byte)

// ConnectionClosed - channel used to signal a connection closing
var ConnectionClosed = make(chan int)

// Register - registers new member & connection
func Register(w http.ResponseWriter, r *http.Request) {

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
    receive: &IncomingMessages,
    closed:  &ConnectionClosed,
  })

  if regErr != nil {
    http.Error(w, regErr.Error(), 500)
    return
  }

  m.Connection.Listen()
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

  // listens for i/o from ws
  go connectionLoop(IncomingMessages, ConnectionClosed, GetAllMembers())

  err := http.ListenAndServe(":3334", r)

  if err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

func connectionLoop(message chan []byte, ConnectionClosed chan int, members map[int]*Member) {
  for {
    select {
    case broadcast := <-message:
      // message received; broadcast to all connections
      log.Println("OUTPUT", string(broadcast))
      for _, member := range members {
        member.Connection.send <- broadcast
      }
    case id := <-ConnectionClosed:
      // connection is signaling close; kill & delete connection
      log.Println("CLOSING", id)
      members[id].Connection.Kill()
      delete(members, id)
    }
  }
}
