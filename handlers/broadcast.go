package handlers

import (
	"fmt"
	"strings"
)

func BroadcastMessages() {
	for msg := range msgs {
		mutex.Lock()
		for _, client := range clients {
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
		mutex.Unlock()
	}
}
