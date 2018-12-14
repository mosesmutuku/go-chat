package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type chatMessage struct {
	sender string
	text   string
}

type chatUser struct {
	name    string
	file    *os.File
	channel chan chatMessage
}

var masterChan = make(chan chatMessage)
var users = make(map[string]*chatUser)

var ctx = context.Background()
var cancel context.CancelFunc

func main() {
	ctx, cancel = context.WithCancel(ctx)

	r := bufio.NewReader(os.Stdin)
	fmt.Println("List the users who will take part in the chat (e.g. user1 user2 user3):")
	userInput, err := r.ReadString('\n')

	if err != nil {
		log.Fatalln("Could not read your response")
	}

	for _, v := range strings.Fields(userInput) {
		users[v] = &chatUser{name: v}
		cb, err := users[v].initUser()

		if err != nil {
			log.Fatalf("There was an error setting up the user: %v \n", users[v].name)
		}

		defer cb()

		go users[v].chat(masterChan)
	}

	go masterChat()

	for ctx.Err() == nil {
		time.Sleep(time.Second * 6)
	}

	close(masterChan)

	for _, u := range users {
		writeLatest(u.file, "\nChat has ended")
	}
}

func (u *chatUser) initUser() (func(), error) {
	f, err := os.Create(u.name + ".txt")
	if err != nil {
		log.Fatalln("Could not create file for user1")
	}

	callBack := func() {
		f.Close()
	}

	u.file = f
	u.channel = make(chan chatMessage)

	return callBack, err
}

func (u *chatUser) chat(masterChan chan<- chatMessage) {
	for ctx.Err() == nil {
		select {
		case message := <-u.channel:
			_, err := writeLatest(u.file, "\n"+message.sender+": "+message.text)

			if err != nil {
				cancel()
			}

			log.Println(u.name, "received message from the channel")
		default:
			// read file
			text, err := readLatest(u.file)

			if err != nil {
				log.Printf("ERROR: There was an error reading the file. User: %+v Error: %v\n", u.name, err)
				cancel()
				return
			}

			if len(text) > 0 {

				if strings.TrimSpace(strings.ToLower(text)) == "quit" {
					log.Printf("User: %v has quit the chat\n", u.name)
					cancel()
					return
				}

				message := chatMessage{
					sender: u.name,
					text:   strings.TrimLeft(text, "\n"),
				}

				masterChan <- message
			}

			time.Sleep(time.Second)
			log.Println(u.name, "after sleeping")
		}
	}
}

func readLatest(f *os.File) (string, error) {
	var bs []byte
	var err error

	for {
		b := make([]byte, 30, 30)
		_, err = f.Read(b)

		if err == nil {
			bs = append(bs, b...)
		} else {
			break
		}
	}

	if err != io.EOF {
		return "", err
	}

	bs = bytes.TrimRight(bs, "\x00")

	log.Println("Read from file", bs)

	return string(bs), nil
}

func writeLatest(f *os.File, s string) (int, error) {
	return f.Write([]byte(s))
}

func masterChat() {
	for ctx.Err() == nil {
		select {
		case message := <-masterChan:
			for k, u := range users {
				if k != message.sender {
					u.channel <- message
				}
			}
		}
	}
}
