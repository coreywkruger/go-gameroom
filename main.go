package main

import (
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

var gotMessage = make(chan []byte)
var isClosing = make(chan int)

func main() {

  // Main http rounter
  r := mux.NewRouter()
  r.Headers("Access-Control-Allow-Origin", "*")
  r.Headers("Access-Control-Allow-Headers", "Content-Type")
  r.Headers("Content-Type", "application/json")
  r.Headers("Content-Type", "text/plain")
  r.HandleFunc("/register", CreateRegistrationHandler(gotMessage, isClosing))

  // listens for i/o from ws
  go connectionLoop(gotMessage, isClosing, GetAllMembers())

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
      log.Println("Broadcasting", string(broadcast))
      for _, member := range members {
        member.Connection.send <- broadcast
      }
    case id := <-ConnectionClosed:
      // connection is signaling close; kill & delete connection
      log.Println("Closing", id)
      members[id].Connection.Kill()
      delete(members, id)
    }
  }
}
