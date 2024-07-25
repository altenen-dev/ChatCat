package handlers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func promptUserName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(clientName, "\n", ""), err

}

// This handler welcomes the new clients and push the old messages to him.
func NewClientsHandler(conn net.Conn) (Client, error) {
	//Print asciiart
	asciiart, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("Error printing asciiArt: ", err)
	}
	coloredAscii := Blue + string(asciiart) + ResetColor
	conn.Write([]byte(coloredAscii))
	//prompt username
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	clientName, err := promptUserName(conn) //promt username
	if err != nil {
		return Client{}, err
	}
	if clientName == "" {
		conn.Write([]byte("Empty name is not allowed!\n"))
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		clientName, err = promptUserName(conn) //promt username
		if err != nil {
			return Client{}, err
		}
	}

	if isClientExist(clientName) {
		conn.Write([]byte("There is a user with this name, please use another name.\n"))
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		clientName, err = promptUserName(conn) //promt username
		if err != nil {
			return Client{}, err
		}
	}
	Mutex.Lock()
	//push msg history
	for _, v := range History {
		conn.Write([]byte(fmt.Sprintln(v)))
	}
	Mutex.Unlock()
	//create the client
	client := Client{Name: clientName, Connection: conn}
	Mutex.Lock()
	//add the client to the map
	Clients[conn] = client
	Mutex.Unlock()
	msg := Message{Content: string(Green) + clientName + " has joined" + ResetColor, sender: conn}
	Msgs <- msg
	Mutex.Lock()
	History = append(History, msg.Content)
	Mutex.Unlock()
	return client, nil
}
