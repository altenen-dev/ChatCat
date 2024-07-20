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

	Mutex.Lock()
	for _, v := range History {
		conn.Write([]byte(fmt.Sprintln(v)))
	}
	Mutex.Unlock()

	clientName = strings.TrimSpace(clientName)
	client := Client{Name: clientName, Connection: conn}

	Mutex.Lock()
	Clients[conn] = client
	Mutex.Unlock()
	msg := Message{Content: clientName + " has joined", sender: conn}
	Msgs <- msg
	History = append(History, msg.Content)
	return client, nil
}
