package main

import (
	"fmt"
	"log"
	"net"
)

const Bufsize int = 256

// enum for types of messages
type MessageType int

const (
	NewConnection MessageType = iota
	Disconnect    MessageType = iota
	Send          MessageType = iota
)

type Message struct {
	Type   MessageType
	Sender Client
	Text   string
}

type Client struct {
	Buffer   []byte
	Username string // Not unique, TODO: make usernames unique (would also have to change DB)
	Conn     net.Conn
}

type Server struct {
	Clients  map[string]*Client
	Messages chan Message
	Listener net.Listener
}

// Handles initial connection for newly registered client, if establishment of connection
// fails, close the connection and ignore
func Connect(conn net.Conn, server Server) {
	client := Client{
		Buffer: make([]byte, Bufsize),
		Conn:   conn,
	}

	_, err := client.Conn.Write([]byte("SERVER_INFO: Welcome to go-chat! Please enter a username: "))
	if err != nil {
		log.Printf("Could not write to: %s\n", conn.RemoteAddr())
		conn.Close()
		return
	}

	n, err := conn.Read(client.Buffer)
	if err != nil {
		log.Printf("Could not read username from: %s\n", conn.RemoteAddr())
		conn.Close()
		return
	}

	client.Username = string(client.Buffer[0 : n-1])

	// Notify other clients of new connection
	server.Messages <- Message{
		Type:   NewConnection,
		Sender: client,
	}

	// Start event loop for client
	go client.Run(server)
}

// Handles event loop for client
func (client Client) Run(server Server) {
	for {
		n, err := client.Conn.Read(client.Buffer)
		if err != nil {
			server.Messages <- Message{
				Type:   Disconnect,
				Sender: client,
			}
			return
		}

		server.Messages <- Message{
			Type:   Send,
			Sender: client,
			Text:   string(client.Buffer[0 : n-1]),
		}
	}
}

func Start(conf GochatConfig) (*Server, error) {
	server := &Server{
		Clients:  map[string]*Client{},
		Messages: make(chan Message),
	}

	// Start listening for tcp connections at `Conf.Port`
	listener, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		return nil, err
	}

	server.Listener = listener

	return server, nil
}

func (server Server) Accept() {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s: %s\n", conn, err)
			continue
		}

		go Connect(conn, server)
	}
}

// TODO: normalize format for sending messages to cchat, json?
// Handles event loop for the server managing forwarding of messages to clients
func (server Server) Run() {
	go server.Accept()

	for {
		message := <-server.Messages
		switch message.Type {

		// New client connected
		case NewConnection:
			server.Clients[message.Sender.Conn.RemoteAddr().String()] = &message.Sender

			outstr := fmt.Sprintf("CONNECT: New user joined with username '%s'\n", message.Sender.Username)
			log.Printf(outstr)
			for _, client := range server.Clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		// Client has disconnected
		case Disconnect:
			message.Sender.Conn.Close()
			delete(server.Clients, message.Sender.Conn.RemoteAddr().String())

			outstr := fmt.Sprintf("DISCONNECT: User '%s@%s' has disconnected\n", message.Sender.Username, message.Sender.Conn.RemoteAddr())
			log.Printf(outstr)
			for _, client := range server.Clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		// Client has sent a message
		case Send:
			outstr := fmt.Sprintf("SEND: %s: %s\n", message.Sender.Username, message.Text)
			log.Printf(outstr)
			for _, client := range server.Clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		}
	}
}
