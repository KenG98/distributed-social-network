package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	cmd := os.Args[1]
	sendIp := "localhost:8080"

	if cmd == "send" {

		req, err := http.NewRequest("GET", sendIp, nil)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		q := req.URL.Query()
		q.Add("name", os.Args[2])
		req.URL.RawQuery = q.Encode()
		client := &http.Client{}
		client.Do(req)
		// resp, err := client.Do(req)
		// log.Println(resp.Body)

	}
	if cmd == "recieve" {

		http.HandleFunc("/", handler)
		http.ListenAndServe(":8080", nil)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query()["name"]
	if len(name) > 0 {
		log.Println("Hello from:", name[0])
		fmt.Fprintf(w, "Hello, "+name[0]+"\n")
	}
	return
}
