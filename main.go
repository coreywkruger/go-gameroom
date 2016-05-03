package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Main http rounter
	r := mux.NewRouter()
	r.Headers("Access-Control-Allow-Origin", "*")
	r.Headers("Access-Control-Allow-Headers", "Content-Type")
	r.Headers("Content-Type", "application/json")
	r.Headers("Content-Type", "text/plain")
	r.HandleFunc("/register", RegistrationHandler)

	// listens for i/o from ws
	StartBroadcast()

	err := http.ListenAndServe(":3334", r)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
