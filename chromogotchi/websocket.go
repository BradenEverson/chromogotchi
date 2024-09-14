package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
		}

	}
}
