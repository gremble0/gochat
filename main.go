package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
)

// IDEAS:
// - Persistent data storage and logging with postgresql or something similar
// - Host multiple chats at once
// - User authentication with usernames and passwords

type MessageType int

const (
	Connect = iota
	Disconnect = iota
	Send = iota
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
	conn.Write([]byte("SERVER_INFO: Welcome to go-chat! Please enter a username: "))

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

// TODO: normalize format for sending messages to cchat, json?
func server(messages chan Message) {
	clients := map[string]*Client{}
	for {
		message := <-messages
		switch message.Type {

		case Connect:
			clients[message.Sender.Conn.RemoteAddr().String()] = &message.Sender

			outstr := fmt.Sprintf("CONNECT: New user joined with username '%s'\n", message.Sender.Username)
			log.Printf(outstr)
			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		case Disconnect:
			message.Sender.Conn.Close()
			delete(clients, message.Sender.Conn.RemoteAddr().String())

			outstr := fmt.Sprintf("DISCONNECT: User '%s@%s' has disconnected\n", message.Sender.Username, message.Sender.Conn.RemoteAddr())
			log.Printf(outstr)
			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		case Send:
			outstr := fmt.Sprintf("SEND: %s: %s\n", message.Sender.Username, message.Text)
			log.Printf(outstr)
			for _, client := range clients {
				if client.Conn.RemoteAddr().String() != message.Sender.Conn.RemoteAddr().String() {
					go client.Conn.Write([]byte(outstr))
				}
			}

		}
	}
}

func dbConnect(host string, port int, user string, password string, dbname string, sslmode string) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to database\n")
	}
	defer db.Close()
	
	rows, _ := db.Query("SELECT * FROM users")
	for rows.Next() {

	}
}

func main() {
	Port := "8080" // Default port unless set by cmd args
	psqlUsername := "postgres" // Default psql user unless set by cmd args
	if len(os.Args) >= 2 {
        psqlUsername = os.Args[1]

		if len(os.Args) >= 3 {
			Port = os.Args[2]
		}
	}

	dbConnect("localhost", 5432, psqlUsername, "", "gochat", "disable")
		
	// Start listening for tcp connections at `Port`
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
