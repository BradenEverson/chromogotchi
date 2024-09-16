package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleNewPetRequest(w http.ResponseWriter, r *http.Request) {
	newPetId := uuid.New().String()
	newPet := makePet()

	allPets[newPetId] = newPet

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"PetId": newPetId,
	}

	json.NewEncoder(w).Encode(response)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	go hanndlConnection(conn)
}

func hanndlConnection(conn *websocket.Conn) {
	var id string
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			request, err := deserializeRequestObject(message)
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			if request.RequestType == "Establish" {
				id = string(request.Metadata)
			}

			if id == "" {
				continue
			}

			_, exists := allPets[id]

			if !exists {
				newPet := makePet()
				allPets[id] = newPet
			}

			pet := allPets[id]

			var responseType string
			var responseData []byte

			var buf bytes.Buffer

			switch request.RequestType {
			case "Feed":
				responseType = "Fed"
				err := binary.Write(&buf, binary.LittleEndian, pet.hunger)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				responseData = buf.Bytes()
			case "Sleep":
				responseType = "Slept"
				err := binary.Write(&buf, binary.LittleEndian, pet.wakefullness)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				responseData = buf.Bytes()
			case "Play":
				responseType = "Happy"
				err := binary.Write(&buf, binary.LittleEndian, pet.depression)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				responseData = buf.Bytes()
			case "Get":
				responseType = "Pet"
				responseData = []byte(pet.name)
			case "Sprite":

				responseType = "Sprite"
				responseData = pet.getSprite()
			}

			if responseType == "" {
				continue
			}

			response := makeRequestObject(responseType, responseData)

			newMessageType := websocket.TextMessage
			newMessage, err := serializeRequestObject(&response)

			if err == nil {
				err = conn.WriteMessage(newMessageType, []byte(*newMessage))
				if err != nil {
					log.Println("Error writing message:", err)
					break
				}
			} else {
				log.Println("Error writing message:", err)
			}
		case websocket.CloseMessage:
			log.Println("Gracefully closing websocket connection")
			break
		}
	}
	conn.Close()
}
