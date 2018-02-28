package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	// make a new DB and add a bunch of messages to it
	db := CreateNewChatDB()
	db.AddMessage("intros", "Ken", "Hi all, I'm Ken.")
	db.AddMessage("intros", "Robert", "Hi Ken, I'm Rob.")
	db.AddMessage("other", "Guy", "Where am I?")
	db.AddMessage("intros", "Samuel", "Hi both of you.")
	db.AddMessage("drones", "Daniel", "Let's talk about DRONES here!!!")
	db.AddMessage("intros", "Ken", "What else can this social network do?")
	db.AddMessage("drones", "Ken", "Oh, I want to talk about drones.")

	// save it to disk
	db.SaveToDisk("saved_network.txt")

	// read it from disk
	db2, _ := LoadFromDisk("saved_network.txt")

	// add some more messages
	db2.AddMessage("drones", "Dude", "I'm here, too, after a save.")

	// print all messages
	db2.PrintDB()
}

type Message struct {
	Sender  string
	Message string
}

type Chatroom struct {
	Name     string
	Messages []Message
}

type ChatDB struct {
	Rooms []Chatroom
}

func (db *ChatDB) SaveToDisk(fn string) error {
	file, err := os.Create(fn)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(*db)
		// for _, r := range db.rooms {
		// 	encoder.Encode(r)
		// 	for _, m := range r.messages {
		// 		encoder.Encode(m)
		// 	}
		// }
		// encoder.Encode(db.rooms)
	}
	file.Close()
	return err
}

func LoadFromDisk(fn string) (*ChatDB, error) {
	db := new(ChatDB)
	file, err := os.Open(fn)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(db)
	}
	file.Close()
	return db, err
}

func CreateNewChatDB() *ChatDB {
	db := new(ChatDB)
	db.Rooms = make([]Chatroom, 0)
	return db
}

func (db *ChatDB) AddMessage(room, sender, msg string) {
	// make the message
	m := &Message{sender, msg}
	// check if the chatroom exists, do your thing
	for i, r := range db.Rooms {
		if r.Name == room {
			db.Rooms[i].Messages = append(r.Messages, *m)
			return
		}
	}
	// the room doesn't exist
	cr := &Chatroom{room, make([]Message, 1)}
	cr.Messages[0] = *m
	db.Rooms = append(db.Rooms, *cr)
	return
}

func (db *ChatDB) PrintDB() {
	for _, cr := range db.Rooms {
		fmt.Println(cr.Name, ":")
		for _, msg := range cr.Messages {
			fmt.Println("\t", msg.Sender, ":", msg.Message)
		}
	}
}
