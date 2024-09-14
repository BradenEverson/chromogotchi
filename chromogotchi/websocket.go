package main

import (
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
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            break
        }

        switch messageType {
        case websocket.TextMessage:
            // Deserialize into a request object

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

