package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"grouper"
	"net/http"
	"os"
	"strconv"
)

type Cli struct {
	gr grouper.Grouper
}

var cli Cli

// this main function is the main function of the entire program.
// It all coordinates here.
func main() {
	cli = Cli{}
	cmd := os.Args[1]
	var port int
	if cmd == "start" {
		port, _ = strconv.Atoi(os.Args[3])
		cli.gr.StartNetwork(os.Args[2], os.Args[3], os.Args[4])
	} else if cmd == "join" {
		port, _ = strconv.Atoi(os.Args[5])
		cli.gr.JoinNetwork(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	}

	serverMuxGrouper := http.NewServeMux()
	serverMuxGrouper.HandleFunc("/handleNewMsg", cli.handleNewMsg)
	go func() {
		http.ListenAndServe(":"+strconv.Itoa(port+1), serverMuxGrouper)
	}()

	fmt.Println("\nWelcome! Type \"quit\" to quit!")
	reader := bufio.NewReader(os.Stdin)
	for {
		entered, _ := reader.ReadString('\n')
		entered = entered[:len(entered)-1]
		if entered == "quit" {
			break
		} else {
			fmt.Println(cli.gr.Me.Name + ": " + entered)
			cli.sendMsgs(cli.gr.Me.Name + ": " + entered)
		}
	}
	fmt.Println("Goodbye.")
}

func (cl *Cli) sendMsgs(msg string) {
	for _, friend := range cli.gr.Them {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(msg)
		friendPortInt, _ := strconv.Atoi(friend.Port)
		friendPort := strconv.Itoa(friendPortInt + 1)
		http.Post("http://"+friend.Ip+":"+friendPort+"/handleNewMsg",
			"application/json; charset=utf-8", b)
	}
}

func (cl *Cli) handleNewMsg(w http.ResponseWriter, r *http.Request) {
	var msg string
	json.NewDecoder(r.Body).Decode(&msg)
	fmt.Println(msg)

}
