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
	Connect    = iota
	Disconnect = iota
	Send       = iota
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

type GochatConfig struct {
	Port   string
	DBConf DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string // NB: Not particularly safe
	DBName   string
}

func dbConnect(DBConf DBConfig) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBConf.Host, DBConf.Port, DBConf.User, DBConf.Password, DBConf.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to database\n")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// rows, _ := db.Query("SELECT * FROM users")
	// for rows.Next() {

	// }
}

func usage() {
	fmt.Printf("%s <gochat-port> <psql-host> <psql-port> <psql-user> <psql-password> <psql-dbname>\n", os.Args[0])
}

func parseConfig(args []string) GochatConfig {
	// Set some defaults for the config
	ret := GochatConfig{
		"8080",
		DBConfig{
			"localhost",
			"5432",
			"postgres",
			"",
			"gochat",
		},
	}

	for i := 1; i < len(args); i++ {
		if i == len(args)-1 {
			log.Printf("%s\n", args[i])
			log.Fatal("Provided flag without argument\n")
		}
		switch args[i] {
		case "-gp":
			ret.Port = args[i+1]
		case "-h":
			ret.DBConf.Host = args[i+1]
		case "-u":
			ret.DBConf.User = args[i+1]
		case "-sp":
			ret.DBConf.Port = args[i+1]
		case "-w":
			ret.DBConf.Password = args[i+1]
		case "-n":
			ret.DBConf.DBName = args[i+1]
		default:
			log.Fatalf("Provided unknown flag '%s'\n", args[i])
		}
	}

	return ret
}

func main() {
	Conf := parseConfig(os.Args)

	dbConnect(Conf.DBConf)

	// Start listening for tcp connections at `Port`
	ln, err := net.Listen("tcp", ":"+Conf.Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s\n", Conf.Port)
	}
	log.Printf("go-chat initialized on port %s\n", Conf.Port)

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
