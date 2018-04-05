package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
      fmt.Println(err)
      return
  }

		for {
			// Read message from browser
      messageType, p, err := conn.ReadMessage()
      if err != nil {
          fmt.Println(err)
          return
      }

      if err := conn.WriteMessage(messageType, p); err != nil {
          fmt.Println(err)
          return
      }
			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(p))

			// Write message back to browser
			if err = conn.WriteMessage(messageType, p); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websocket.html")
	})

	http.ListenAndServe(":8080", nil)
}
