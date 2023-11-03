package main

import (
	"fmt"
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
	Type   MessageType
	Sender Client
	Text   string
}

type Client struct {
	Username string
	Conn     net.Conn
}

func client(conn net.Conn, messages chan Message) {
	buf := make([]byte, 256)
	conn.Write([]byte("Welcome to go-chat! Please enter a username: "))

	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Could not read username from: %s\n", conn.RemoteAddr())
		conn.Close()
		return
	}

	client := Client{
		Username: string(buf[0 : n-1]),
		Conn:     conn,
	}

	messages <- Message{
		Type:   Connect,
		Sender: client,
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			messages <- Message{
				Type:   Disconnect,
				Sender: client,
			}
			return
		}

		messages <- Message{
			Type:   Send,
			Sender: client,
			Text:   string(buf[0 : n-1]),
		}
	}
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	for {
		message := <-messages
		switch message.Type {

		case Connect:
			clients[message.Sender.Conn.RemoteAddr().String()] = &message.Sender

			outstr := "New user joined with username '" + message.Sender.Username + "'\n"
			log.Printf(outstr)

			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					client.Conn.Write([]byte(outstr))
				}
			}

		case Disconnect:
			message.Sender.Conn.Close()
			delete(clients, message.Sender.Conn.RemoteAddr().String())

			outstr := fmt.Sprintf("User '%s@%s' has disconnected\n", message.Sender.Username, message.Sender.Conn.RemoteAddr())
			log.Printf(outstr)
			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					client.Conn.Write([]byte(outstr))
				}
			}

		case Send:
			outStr := message.Sender.Username + ": " + message.Text + "\n"
			log.Printf(outStr)

			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					client.Conn.Write([]byte(outStr))
				}
			}

		}
	}
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
