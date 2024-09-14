package main

import (
	"fmt"
	"log"
	"net/http"
)

var allPets map[string]Pet
func main() {
	http.HandleFunc("/connection", handleWebSocket)
	http.HandleFunc("/newpet", handleNewPetRequest)

	fmt.Println("WebSocket server started on :7878")
	err := http.ListenAndServe(":7878", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
