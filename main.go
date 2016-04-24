package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// BroadcastCh -
var BroadcastCh = make(chan Message)

// ClosingCh -
var ClosingCh = make(chan int)

func main() {

	// Main http rounter
	r := mux.NewRouter()
	r.Headers("Access-Control-Allow-Origin", "*")
	r.Headers("Access-Control-Allow-Headers", "Content-Type")
	r.Headers("Content-Type", "application/json")
	r.Headers("Content-Type", "text/plain")
	r.HandleFunc("/register", RegistrationHandler(&BroadcastCh))

	// listens for i/o from ws
	go func() {
		for {
			select {
			case broadcast := <-BroadcastCh:
				// BroadcastCh received; broadcast to all connections
				connections, _ := GetAllConnections()
				log.Println("Broadcasting", string(broadcast.message))
				for _, connection := range connections {
					if connection.ID != broadcast.ID {
						*connection.send <- broadcast.message
					}
				}
			}
		}
	}()

	err := http.ListenAndServe(":3334", r)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
