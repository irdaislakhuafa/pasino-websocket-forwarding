package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	port = "8080"
	url  = "wss://socket.pasino.com/dice/"
)

func main() {
	ctx := context.Background()
	webSocketDialler := websocket.DefaultDialer
	headers := http.Header{}

	webSocketConnection, _, err := webSocketDialler.DialContext(ctx, url, headers)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer webSocketConnection.Close()

	webSockerUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// logic to allow origin here
			return true
		},
	}

	server := http.DefaultServeMux
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		clientConnection, err := webSockerUpgrader.Upgrade(w, r, nil)
		if err != nil {
			message := fmt.Sprintf("Failed to upgrade HTTP connection to WebSockets: %v", err)
			log.Println(message)
			ReturnErrorResponse(w, message)
			return
		}
		defer clientConnection.Close()

		// Forwarding data from client to web socket server
		for {
			messageType, data, err := webSocketConnection.ReadMessage()
			if err != nil {
				message := fmt.Sprintf("Error reading from client: %v", err)
				log.Println(message)
				ReturnErrorResponse(w, message)
				return
			}

			// is message is close message
			if messageType == websocket.CloseMessage {
				message := "Receive close message"
				log.Println(message)

				err := webSocketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					message := fmt.Sprintf("Error writing close message: %v", err)
					log.Println(message)
					ReturnErrorResponse(w, message)
				}
				return
			} else {
				err = webSocketConnection.WriteMessage(messageType, data)
				if err != nil {
					message := fmt.Sprintf("Error writing to websocket server: %v", err)
					log.Println(message)
				}
				return
			}
		}
	})

	var handler http.Handler = server

	log.Printf("Server listening on port %v\n", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalln("Failed to start HTTP server:", err)
	}
}

func ToStringJSON[T any](value T) string {
	jsonBytes, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		message := fmt.Sprintf("Error marshalling JSON: %+v", err)
		log.Println(message)
		return fmt.Sprintf(`{"error": "%+v"}`, message)
	}
	return string(jsonBytes)
}

func ReturnErrorResponse[T any](w http.ResponseWriter, message T) {
	w.Write([]byte(
		ToStringJSON(
			map[string]interface{}{
				"error": fmt.Sprintf("%+v", message),
			},
		),
	))
}
