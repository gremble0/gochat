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

func server(messages chan Message) {
	clients := map[string]*Client{}
	for {
		message := <- messages
		switch message.Type {
		case Connect:
			conn := message.Sender.Conn
			usernameBuf := make([]byte, 20)
			conn.Write([]byte("Welcome to go-chat! Please enter a username:\n"))
			
			n, err := message.Sender.Conn.Read(usernameBuf)
			if err != nil {
				log.Printf("Could not read username from: %s\n", conn.RemoteAddr())
				message.Sender.Conn.Close()
				messages <- Message {
					Type: Disconnect,
					Sender: message.Sender,
				}
			}

			clients[message.Sender.Conn].Username = string(usernameBuf[0:n])
		case Disconnect:
		case Send:
		}
	}
}

func handleConnection(conn net.Conn, messages chan Message) {
	buf := make([]byte, 512)

	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
	// conn.Write([]byte("Welcome to go chat! Please enter a username:\n"))

	// n, err := conn.Read(buf)

	client := Client {
		Username: string(buf[0:n]),
		Conn: conn,
	}

	for {
		conn.Read(buf)
		messages <- Message {
			Type: Send,
			Sender: client,
			Text: string(buf[0:n]),
		}
		// log.Printf("%s: %s\n", client.Username, buf)
		conn.Write(buf[0:n])
	}
}

func main() {
	ln, err := net.Listen("tcp", ":" + Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s\n", Port)
	}

	messages := make(chan Message)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s\n", conn)
			continue
		}

		messages <- Message {
			Type: Connect,
		}
		go handleConnection(conn, messages)
	}
}
