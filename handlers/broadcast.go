package handlers

import (
	"fmt"
	"strings"
)

func BroadcastMessages() {
	for msg := range Msgs {
		Mutex.Lock()
		for _, client := range Clients {
			if msg.sender != client.Connection && (strings.Contains(msg.Content, "left") || strings.Contains(msg.Content, "join")) {
				prompt := FormatPrompt(client)
				client.Connection.Write([]byte("\n"))
				client.Connection.Write([]byte(msg.Content))
				client.Connection.Write([]byte("\n"))
				client.Connection.Write([]byte(prompt))
			}
			if msg.sender != client.Connection {
				client.Connection.Write([]byte(fmt.Sprintln(msg.Content)))
			}
		}
		Mutex.Unlock()
	}
}
