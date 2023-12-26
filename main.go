package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// creating and configuring an upgrader object
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Client holds info about connection
type Client struct {
	ID   string          // A unique identifier
	conn *websocket.Conn // A pointer to a websocket.Conn object
	send chan []byte     // A channel of byte slices
}

// Generate a unique ID for each client
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Define a map to store connected clients
var clients = make(map[*Client]bool)

// Read messages from the WebSocket connection
func (client *Client) readPump() {
	defer func() {
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		// add client ID to message
		message = []byte(fmt.Sprintf("%s: %s", client.ID, string(message)))

		// send the message to every connected client
		for client := range clients {
			client.send <- message
		}
	}
}

// Write message to the WebSocket connection
func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()

	for {
		message, ok := <-client.send

		if !ok {
			// The channel was closed
			break
		}

		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

// Handle new WebSocket connection
func handleConnections(writer http.ResponseWriter, request *http.Request) {
	// Upgrade initial GET request to a WebSocket
	websocket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		return
	}

	client := &Client{
		ID:   generateUniqueID(),
		conn: websocket,
		send: make(chan []byte),
	}

	clients[client] = true

	// Start listening for incoming chat messages
	go client.readPump()

	// Start sending messages obtained from the broadcast channel
	go client.writePump()

	// Send message "user ID connected" to every connected client except the new one
	for c := range clients {
		if c != client {
			c.send <- []byte("user " + client.ID + " connected")
		}
	}
}

func main() {
	// Serve files from the current directory
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "index.html")
	})
	http.HandleFunc("/script.js", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "script.js")
	})

	// Handle WebSocket requests
	http.HandleFunc("/ws", handleConnections)

	// Start the server
	fmt.Println("Server started on localhost:8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Printf("Server failed to start: %s\n", err)
	}
}
