package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

// IDEAS:
// - Host multiple chats at once
// - User authentication with usernames and passwords

type GochatConfig struct {
	Port   string
	DBConf DBConfig
}

func usage() {
	fmt.Printf("Usage: %s [OPTION]...\n", os.Args[0])
	fmt.Printf("    --help               display this help message\n")
	fmt.Printf("    -gp <port>          set the port for gochat to listen to connections on\n")
	fmt.Printf("    -hn <hostname>      set the hostname for the postgres database connection\n")
	fmt.Printf("    -u  <username>      set the username for the postgres database connection\n")
	fmt.Printf("    -sp <port>          set the port for the postgres database connection\n")
	fmt.Printf("    -w  <password>      set the password for the postgres database connection (NB: NOT SECURE)\n")
	fmt.Printf("    -n  <database name> set the database name for the postgres database connection\n")

	os.Exit(1)
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

	for i := 1; i < len(args); i = i + 2 {
		if i == len(args) - 1 {
			log.Fatalf("Provided flag without argument\n")
		}

		switch args[i] {
		case "--help":
			usage()
		case "-gp":
			ret.Port = args[i+1]
		case "-hn":
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

	// Initialize database
	_, err := dbConnect(Conf.DBConf)
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	log.Printf("Successfully connected to the '%s' database\n", Conf.DBConf.DBName)

	// Start listening for tcp connections at `Conf.Port`
	ln, err := net.Listen("tcp", ":"+Conf.Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", Conf.Port, err)
	}
	log.Printf("go-chat initialized on port %s\n", Conf.Port)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s: %s\n", conn, err)
			continue
		}

		go client(conn, messages)
	}
}
