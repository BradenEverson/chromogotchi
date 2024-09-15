package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const SLEEP_TIME = 120

var allPets map[string]Pet = make(map[string]Pet)

func updatePetAttributes() {
	for {
		time.Sleep(SLEEP_TIME * time.Second)
		fmt.Println("Updating...")
		for _, pet := range allPets {
			pet.hunger -= pet.hunger
			pet.happiness -= pet.happiness
			pet.wakefullness -= pet.wakefullness
		}
	}
}

func main() {
	http.HandleFunc("/connection", handleWebSocket)
	http.HandleFunc("/newpet", handleNewPetRequest)

	fmt.Println("WebSocket server started on :7878")

	go updatePetAttributes()
	err := http.ListenAndServe(":7878", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
