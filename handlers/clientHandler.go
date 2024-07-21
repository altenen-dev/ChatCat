package handlers

import (
	"bufio"
)

// This handler reads the messages from the client to send it to the channel, and adds it to the history
func HandleClient(client Client) {
	// #note that this hanlder will be running as long  as the client is in the server
	//When the client leavs, this func will terminate, and before that this defer func will execute:
	defer func() {
		//client exit
		Mutex.Lock()
		delete(Clients, client.Connection)
		Mutex.Unlock()
		msg := Message{Content: client.Name + " has left", sender: client.Connection}
		Msgs <- msg
		History = append(History, msg.Content)
		client.Connection.Close()
	}()
	for {
		//Read the msg form the user
		message, err := bufio.NewReader(client.Connection).ReadString('\n')
		if err != nil {
			return
		}
		// skipping the empty messages
		if message == "" {
			client.Connection.Write([]byte(FormatPrompt(client)))
			continue
		}
		fMsg := FormatMessage(client, message)
		msg := Message{Content: fMsg, sender: client.Connection}
		Msgs <- msg
		History = append(History, msg.Content)

	}
}
