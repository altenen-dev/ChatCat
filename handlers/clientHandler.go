package handlers

import (
	"bufio"
	"fmt"
	"strings"
)

func HandleClient(client Client) {
	defer func() {
		mutex.Lock()
		delete(clients, client.Connection)
		mutex.Unlock()
		msg := Message{Content: client.Name + " has left\n", sender: client.Connection}
		msgs <- msg
		history = append(history, msg.Content)
		client.Connection.Close()
	}()
	promptMsg := FormatPrompt(client)
	for {
		client.Connection.Write([]byte(promptMsg))
		message, err := bufio.NewReader(client.Connection).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input:", err)
			return
		}
		message = strings.TrimSpace(message) // Trim leading/trailing spaces
		if message == "" {
			continue
		}
		fMsg := FormatMessage(client, message)
		msg := Message{Content: fMsg, sender: client.Connection}
		msgs <- msg
		history = append(history, msg.Content)

	}
}
