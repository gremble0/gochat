package main

import (
	"log"
	"net"
)

const (
	Port = "8080"
)

func handleConnection(conn net.Conn) {
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":" + Port)
	if err != nil {
		log.Fatalf("Could not listen to port %s\n", Port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection from %s\n", conn)
		}

		go handleConnection(conn)
	}
}
