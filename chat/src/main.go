package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Data)

// Upgrader takes normal HTTP connection and upgrades it to a WebSocket
var upgrader = websocket.Upgrader{}

// Definition of the data object (formerly Message object)
type Data struct{
  Timestamp string `json:"timestamp"`
  Username string `json:"username"`
  Action string `json:"action"`
  Contents string `json:"contents"`
  // e.g. to send message, action would be 'message' and contents would be the message
  // Other ctions include joining server, leaving server, changing username
  // Add authentication in future but will require another struct
  // Changed: Added timestamp field, removed email field
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
    var incData Data
    // Read in data/message as JSON and map to data object
    err := ws.ReadJSON(&incData)
    if err != nil {
      log.Printf("error: %v", err)
      delete(clients, ws)
      break
    }
    // Debugging log to console feature
    log.Printf("%s: %s\n", incData.Username, incData.Contents)
    // Send newly received message to broadcast channel
    broadcast <- incData
  }
}

func handleMessages() {
  for {
    // Grab next message from broadcast channel
    incData := <-broadcast
    // Add function to save message to database file on server
    // Send out to every client connected
    for client := range clients {
      err := client.WriteJSON(incData)
      if err != nil {
        log.Printf("error: %v", err)
        client.Close()
        delete(clients, client)
      }
    }
  }
}
