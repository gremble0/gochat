package main

import (
	"log"
	"net"
)

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
	Type MessageType
	Sender Client
	Text string
}

type Client struct {
	Username string
	Conn net.Conn
}

func client(client Client, messages chan Message) {
	// buf := make([]byte, 256)
	for {
		// client.Conn.Read(b []byte)
		messages <- handleSend(client)
	}
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	for {
		message := <- messages
		switch message.Type {
		case Connect:
			clients[message.Sender.Conn.RemoteAddr().String()] = &message.Sender
			// message.Sender.Conn.Write([]byte(fmt.Sprintf("New user joined with username: %s", message.Sender.Username)))
			log.Printf("New user joined with username %s", message.Sender.Username)

			go client(message.Sender, messages)
		case Disconnect:
		case Send:
			messages <- handleSend(message.Sender)
		}
	}
}

func handleConnect(conn net.Conn) Client {
	usernameBuf := make([]byte, 20)
	conn.Write([]byte("Welcome to go-chat! Please enter a username:\n"))
	
	n, err := conn.Read(usernameBuf)
	if err != nil {
		log.Printf("Could not read username from: %s\n", conn.RemoteAddr())
		conn.Close()
		return Client {}
	}

	return Client {
		Username: string(usernameBuf[0:n]),
		Conn: conn,
	}
}

func handleDisconnect(client Client) Message {
	log.Printf("User %s @%s has disconnected", client.Username, client.Conn.RemoteAddr())
	client.Conn.Close()

	return Message {
		Type: Disconnect,
		Sender: client,
	}
}

func handleSend(client Client) Message {
	messageBuf := make([]byte, 256)
	conn := client.Conn

	n, err := conn.Read(messageBuf)
	if err != nil {
		return handleDisconnect(client)
	}

	return Message {
		Type: Send,
		Sender: client,
		Text: string(messageBuf[0:n]),
	}
}

func main() {
	ln, err := net.Listen("tcp", ":" + Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s\n", Port)
	}

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s\n", conn)
			continue
		}

		messages <- Message {
			Type: Connect,
			Sender: handleConnect(conn),
		}
	}
}
