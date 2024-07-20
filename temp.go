package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Name       string
	Connection net.Conn
}

type Message struct {
	Content string
	sender  net.Conn
}

var (
	clients = make(map[net.Conn]Client)
	msgs    = make(chan Message) // Initialize the channel
	mutex   sync.Mutex
	history []string
)

func main() {
	port := "8989" // Default port

	dataStream, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TCP server is running on:", port)
	defer dataStream.Close()

	go broadcastMessages() // Start a goroutine to broadcast messages
	for {
		conn, err := dataStream.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		if len(clients) >= 10 {
			conn.Write([]byte("Max connections reached, rejecting new connection\n"))
			conn.Close()
			continue
		}

		conn.Write([]byte("[ENTER YOUR NAME]: "))
		reader := bufio.NewReader(conn)
		clientName, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading client name:", err)
			conn.Close()
			continue
		}

		for _, v := range history {
			conn.Write([]byte(fmt.Sprintln(v)))
		}

		clientName = strings.TrimSpace(clientName)
		client := Client{Name: clientName, Connection: conn}

		mutex.Lock()
		clients[conn] = client
		mutex.Unlock()
		msg := Message{Content: clientName + " has joined\n", sender: conn}
		msgs <- msg
		history = append(history, msg.Content)

		go handleClient(client)
	}
}

func handleClient(client Client) {
	defer func() {
		mutex.Lock()
		delete(clients, client.Connection)
		mutex.Unlock()
		msg := Message{Content: client.Name + " has left\n", sender: client.Connection}
		msgs <- msg
		history = append(history, msg.Content)
		client.Connection.Close()
	}()
	name := strings.TrimSpace(client.Name)

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	template := fmt.Sprintf("[%s][%s]: ", timeStamp, name)
	for {
		client.Connection.Write([]byte(template))
		message, err := bufio.NewReader(client.Connection).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input:", err)
			return
		}
		message = strings.TrimSpace(message) // Trim leading/trailing spaces
		if message == "" {
			continue
		}
		cMsg := fmt.Sprintf("%s%s", template, message)

		msg := Message{Content: cMsg, sender: client.Connection}
		msgs <- msg
		history = append(history, msg.Content)

	}
}

func broadcastMessages() {
	for msg := range msgs {
		mutex.Lock()
		for _, client := range clients {
			if msg.sender != client.Connection && (strings.Contains(msg.Content, "left") || strings.Contains(msg.Content, "join")) {
				timeStamp := time.Now().Format("2006-01-02 15:04:05")
				template := fmt.Sprintf("[%s][%s]: ", timeStamp, client.Name)
				client.Connection.Write([]byte(fmt.Sprintln(msg.Content)))
				client.Connection.Write([]byte(template))
			}
			if msg.sender != client.Connection {
				client.Connection.Write([]byte(fmt.Sprintln(msg.Content)))
			}
		}
		mutex.Unlock()
	}
}
