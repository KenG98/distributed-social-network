package server

import "fmt"

/*
	This is responsible for the local server which interacts with the local user's
	interface. It will need access to the database, and will need to display everything
	the user needs to interact with the social network
*/

func StartLocal() {
	fmt.Println("Starting local server")
}

// open a connection to the database

// create an http server and serve it on a specified port

// serve the page with a list of chatrooms

// serve the page which renders a specific chatroom

// serve the page which lets a user manage their identity
// (they should be able to set a username for the chatrooms)

// use the pages in ui/pages...
