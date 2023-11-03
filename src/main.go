package main

import (
	"log"
	"net"
	"./client"
)

// IDEAS:
// - Client with raylib in C
// - Persistent data storage and logging with postgresql or something similar
// - Host multiple chats at once
// - User authentication with usernames and passwords

const (
	Port = "8080"
)

type MessageType int

const (
	Connect = iota + 1
	Disconnect
	Send
)

type Message struct {
	Type   MessageType
	Sender Client
	Text   string
}

type Client struct {
	Username string
	Conn     net.Conn
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s\n", Port)
	}
	log.Printf("go-chat initialized on port %s\n", Port)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s\n", conn)
			continue
		}

		go client(conn, messages)
	}
}
