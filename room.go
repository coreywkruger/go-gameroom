package main

import (
  "log"
)

// Room -
type Room struct {
  Broacast chan []byte
}

// Start -
func (r *Room) Start(members map[int]*Member) {
  go func() {
    for {
      select {
      case broadcast := <-r.Broacast:
        log.Println("OUTPUT", string(broadcast))
        for _, member := range members {
          member.Connection.send <- broadcast
        }
      }
    }
  }()
}
