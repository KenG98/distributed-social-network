package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if args[1] == "write" {
		saveSentences(args[2])
	} else if args[1] == "read" {

	}
}

func saveSentences(fn string) {
	// ask the user to enter a bunch of strings
	fmt.Println("Enter a sentence and press enter. Enter a blank line to stop.")

	doneEntering := false
	var sentences []string
	reader := bufio.NewReader(os.Stdin)

	for !doneEntering {
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			break
		} else {
			sentences = append(sentences, text)
		}
	}

	// saving to file
	file, _ := os.Create(fn)
	defer file.Close()
	for _, s := range sentences {
		file.WriteString(s)
	}

	// confirm and exit
	fmt.Println("Okay!")
}
