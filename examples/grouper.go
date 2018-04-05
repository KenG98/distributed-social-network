package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

type peer struct {
	Name string
	Ip   string
	Port string
}

type grouper struct {
	mu   sync.Mutex
	Me   peer
	Them []peer
}

var gr grouper

// USAGE

// run 'grouper start ip port name' to begin your own network
// (params are your ip and your name)
// run 'grouper join ipFriend portFriend ipSelf portSelf name' to join a network
// (ipFriend is the node alread in the network, self and name
// for this actual node)

// This function will start an HTTP server no matter what,
// the process will be different though, depending on ^^

func main() {
	cmd := os.Args[1]
	gr = grouper{}

	// handle a start, join, or help command
	if cmd == "start" {
		go listenToSIGINT()
		gr.startNetwork(os.Args[2], os.Args[3], os.Args[4])
	} else if cmd == "join" {
		go listenToSIGINT()
		gr.joinNetwork(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	} else if cmd == "help" {
		fmt.Println("To start a network:")
		fmt.Println("\tgrouper start ip_addr port_number name")
		fmt.Println("To join a network:")
		fmt.Println("\tgrouper join friend_ip friend_port my_ip my_port name")
	} else {
		fmt.Println("Unrecognized command:", cmd)
	}
}

func handleGetPeers(w http.ResponseWriter, r *http.Request) {
	allUsers := append(gr.Them, gr.Me)
	json.NewEncoder(w).Encode(allUsers)
}

func handleJoinNet(w http.ResponseWriter, r *http.Request) {
	user := peer{}
	json.NewDecoder(r.Body).Decode(&user)
	gr.Them = append(gr.Them, user)
	fmt.Println(user.Name, "has joined")
}

func handleLeaveNet(w http.ResponseWriter, r *http.Request) {
	user := peer{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user.Name, "has left")
	// find who it is, and remove them
	for ind, usr := range gr.Them {
		if user == usr {
			gr.Them = append(gr.Them[0:ind], gr.Them[ind+1:len(gr.Them)]...)
			break
		}
	}
}

// func - start your own network
func (gr *grouper) startNetwork(myIp, myPort, myName string) {
	fmt.Println("starting network as", myIp+":"+myPort+":"+myName)
	gr.Me = peer{Ip: myIp, Port: myPort, Name: myName}
	// start serving
	http.HandleFunc("/getPeers", handleGetPeers)
	http.HandleFunc("/joinNet", handleJoinNet)
	http.HandleFunc("/leaveNet", handleLeaveNet)
	http.ListenAndServe(":"+myPort, nil)
}

// func - join a network
func (gr *grouper) joinNetwork(friendIp, friendPort, myIp, myPort, myName string) {
	fmt.Println("joining network as", myIp+":"+myPort+":"+myName)
	gr.Me = peer{Ip: myIp, Port: myPort, Name: myName}
	// get all peers, set them to my peers
	gr.getPeers(peer{Ip: friendIp, Port: friendPort})
	// tell everyone i'm joining
	gr.sendJoinRequests()
	// start serving
	http.HandleFunc("/getPeers", handleGetPeers)
	http.HandleFunc("/joinNet", handleJoinNet)
	http.HandleFunc("/leaveNet", handleLeaveNet)
	http.ListenAndServe(":"+myPort, nil)

}

func (gr *grouper) getPeers(friend peer) {
	cli := &http.Client{}
	r, err := cli.Get("http://" + friend.Ip + ":" + friend.Port + "/getPeers")
	if err != nil {
		return
	}
	defer r.Body.Close()
	var otherUsers []peer
	err = json.NewDecoder(r.Body).Decode(&otherUsers)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	gr.Them = otherUsers
}

func (gr *grouper) sendJoinRequests() {
	for _, usr := range gr.Them {
		go func(p peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(gr.Me)
			http.Post("http://"+p.Ip+":"+p.Port+"/joinNet", "application/json; charset=utf-8", b)
		}(usr)
	}
}

func (gr *grouper) sendLeaveRequests() {
	var wg sync.WaitGroup
	for _, usr := range gr.Them {
		wg.Add(1)
		go func(p peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(gr.Me)
			http.Post("http://"+p.Ip+":"+p.Port+"/leaveNet", "application/json; charset=utf-8", b)
			wg.Done()
		}(usr)
		wg.Wait()
	}
}

func listenToSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Leaving the network.")
		gr.sendLeaveRequests()
		os.Exit(0)
	}()
}
