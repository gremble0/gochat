package main

import (
	"fmt"
	"log"
)

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
