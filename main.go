package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var allPets map[string]Pet

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleNewPetRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Pet Request")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			fmt.Println("Text")
		case websocket.BinaryMessage:
			fmt.Println("Binary")
			if len(message) == 16 {
				petId := string(message)
				petName := string(message[16:])
				allPets[petId] = makePet(petName)
				fmt.Printf("New Pet %s with name %s\n", petId, petName)
			}
		}

		fmt.Printf("Received: %s\n", message)
		newMessageType := websocket.TextMessage
		newMessage := []byte("Hello")

		err = conn.WriteMessage(newMessageType, newMessage)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/newpet", handleNewPetRequest)

	fmt.Println("WebSocket server started on :7878")
	err := http.ListenAndServe(":7878", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
