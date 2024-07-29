package handlers

import (
	"fmt"
)

// This handler prints the messages to all other clients, and print the prompt msg
func BroadcastMessages() {
	for msg := range Msgs {
		Mutex.Lock()
		for _, client := range Clients {
			prompt := FormatPrompt(client)
			if msg.sender != client.Connection {
				_, err := client.Connection.Write([]byte("\n" + fmt.Sprintln(msg.Content) + prompt))
				if err != nil {
					fmt.Println("Error writing the message:", err)
				}
			} else {
				_, err := client.Connection.Write([]byte(prompt))
				if err != nil {
					fmt.Println("Error writing the message:", err)
				}
			}
		}
		Mutex.Unlock()
	}
}
