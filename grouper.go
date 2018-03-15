package main

import (
	"fmt"
)

// USAGE

// run 'grouper start ip name' to begin your own network
// (params are your ip and your name)
// run 'grouper join ipFriend ipSelf name' to join a network
// (ipFriend is the node alread in the network, self and name
// for this actual node)

// This function will start an HTTP server no matter what,
// the process will be different though, depending on ^^

func main() {

}

type grouper struct {
	name  string
	ip    string
	peers []grouper
}

// func - start your own network
func (gr *grouper) startNetwork(myIp, myName string) {

}

// func - join a network
func (gr *grouper) joinNetwork(friendIp, myIp, myName string) {

}

// func - request and add peers from connecting node

// func - return peers to connecting node request

// func - add a peer when the peer requests to join the network

// func - leave the network
