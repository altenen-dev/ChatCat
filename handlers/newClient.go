package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func NewClientsHandler(conn net.Conn) (Client, error) {
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	reader := bufio.NewReader(conn)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		return Client{}, err
	}

	mutex.Lock()
	for _, v := range history {
		conn.Write([]byte(fmt.Sprintln(v)))
	}
	mutex.Unlock()

	clientName = strings.TrimSpace(clientName)
	client := Client{Name: clientName, Connection: conn}

	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()
	msg := Message{Content: clientName + " has joined\n", sender: conn}
	msgs <- msg
	history = append(history, msg.Content)
	return client, nil
}
