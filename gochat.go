package main

import (
	"fmt"
	"log"
	"os"
)

// IDEAS:
// - Send recent messages from database to new connections
// - Host multiple chats at once
// - User authentication with usernames and passwords

// GochatConfig contains some configuration options for the gochat server
type GochatConfig struct {
	Port string
	dbc  DBConfig
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

// parseConfig generates a config from a command line-like string array
func parseConfig(args []string) GochatConfig {
	// Set some defaults for the config
	ret := GochatConfig{
		Port: "8080",
		dbc: DBConfig{
			"localhost",
			"5432",
			"postgres",
			"",
			"gochat",
		},
	}

	for i := 1; i < len(args); i = i + 2 {
		if i == len(args)-1 {
			log.Fatalf("Provided flag without argument\n")
		}

		switch args[i] {
		case "--help":
			usage()
		case "-gp":
			ret.Port = args[i+1]
		case "-hn":
			ret.dbc.Host = args[i+1]
		case "-u":
			ret.dbc.User = args[i+1]
		case "-sp":
			ret.dbc.Port = args[i+1]
		case "-w":
			ret.dbc.Password = args[i+1]
		case "-n":
			ret.dbc.DBName = args[i+1]
		default:
			log.Fatalf("Provided unknown flag '%s'\n", args[i])
		}
	}

	return ret
}

// main hosts the gochat server
func main() {
	conf := parseConfig(os.Args)

	// Initialize database
	db, err := dbConnect(conf.dbc)
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	log.Printf("Successfully connected to the '%s' database\n", conf.dbc.DBName)

	// Initialize server
	server, err := Start(conf, db)
	if err != nil {
		log.Fatalf("Could not listen to port %s: %s\n", conf.Port, err)
	}
	log.Printf("go-chat initialized on port %s\n", conf.Port)

	server.Run()
}
