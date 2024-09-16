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
		for key, pet := range allPets {
			pet.updateHunger(pet.HungerRate * -1)
			pet.updateHappy(pet.Depression * -1)
			pet.updateSleep(pet.SleepyRate * -1)

			allPets[key] = pet
		}
	}
}

func main() {
	err := connectClient()
	if err != nil {
        log.Println("DB Connect Error, defaulting to empty pets:\n", err)
	} else {
        loadPetsFromMongo()
        go savePetsToMongo()
    }

	http.HandleFunc("/connection", handleWebSocket)
	http.HandleFunc("/newpet", handleNewPetRequest)

	fmt.Println("WebSocket server started on :7878")

	go updatePetAttributes()
	err = http.ListenAndServe(":7878", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
