package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Upgrader takes normal HTTP connection and upgrades it to a WebSocket
var upgrader = websocket.Upgrader{}

// Definition of the message object
type Message struct{
  Email string `json:"email"`
  Username string `json:"username"`
  Message string `json:"message"`
}

func main() {
  // Creates the file server
  fs := http.FileServer(http.Dir("../public"))
  http.Handle("/", fs)

  // Configure websocket route
  http.HandleFunc("/ws", handleConnections)

  // Listen for incoming chat messages using a goroutine
  go handleMessages()

  // Start server on localhost:8000 and log errors
  log.Println("HTTP server initialized on 127.0.0.1:8000")
  err := http.ListenAndServe(":8000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
  // Upgrade GET request to websocket
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }

  // Close connection when function returns
  defer ws.Close()

  // Register new client (adds to global clients map)
  clients[ws] = true

  for {
    var msg Message
    // Read in message as JSON and map to message object
    err := ws.ReadJSON(&msg)
    if err != nil {
      log.Printf("error: %v", err)
      delete(clients, ws)
      break
    }
    // Send newly received message to broadcast channel
    broadcast <- msg
  }
}

func handleMessages() {
  for {
    // Grab next message from broadcast channel
    msg := <-broadcast
    // Send out to every client connected
    for client := range clients {
      err := client.WriteJSON(msg)
      if err != nil {
        log.Printf("error: %v", err)
        client.Close()
        delete(clients, client)
      }
    }
  }
}
