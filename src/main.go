package main

import "fmt"
import "server"

// This starts the entire program

func main() {
	fmt.Println("Starting the Distributed Social Network")
	server.StartNetwork()
	server.StartLocal()
}
