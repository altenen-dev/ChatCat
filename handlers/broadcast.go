package handlers

import (
	"fmt"
)

func BroadcastMessages() {
	for msg := range Msgs {
		Mutex.Lock()
		for _, client := range Clients {
			prompt := FormatPrompt(client)
			if msg.sender != client.Connection {
				client.Connection.Write([]byte("\n"))
				client.Connection.Write([]byte(fmt.Sprintln(msg.Content)))
				// client.Connection.Write([]byte("\n"))
				client.Connection.Write([]byte(prompt))
			} else {
				client.Connection.Write([]byte(prompt))
			}
		}
		Mutex.Unlock()
	}
}
