package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	users          []User
	history        []string
	maxConnections = 10
)

func main() {
	port := "8989" // Default port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	dataStream, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TCP server is running on: https://", dataStream.Addr().String())
	defer dataStream.Close()

	for {
		if len(users) >= maxConnections {
			fmt.Println("Max connections reached, rejecting new connection")
			continue
		}

		connection, err := dataStream.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleRequests(connection)
	}
}

func handleRequests(conn net.Conn) {
	defer conn.Close()

	// Read the user's name
	nameBuffer := bufio.NewReader(conn)
	clientName, err := nameBuffer.ReadString('\n')
	if err != nil {
		log.Println("Error reading name from client:", err)
		return
	}
	clientName = strings.TrimSpace(clientName)

	// Check if the user name is valid
	if clientName == "" {
		conn.Write([]byte("Name cannot be empty\n"))
		return
	}

	// Add the user to the list
	user := User{name: clientName, conn: conn}
	users = append(users, user)
	sendMessageToAll(clientName+" has joined", conn)
	// Send the message history to the new user
	for _, msg := range history {
		conn.Write([]byte(msg + "\n"))
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			removeUser(conn)
			sendMessageToAll(clientName+" disconnected", conn)
			return
		}

		clientMsg := strings.TrimSpace(string(buffer[:n]))
		if clientMsg != "" {
			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			msg := fmt.Sprintf("[%s][%s]: %s", timeStamp, clientName, clientMsg)
			history = append(history, msg)
			sendMessageToAll(msg, conn)
		}
	}
}

func sendMessageToAll(msg string, conn net.Conn) {
	for _, user := range users {
		if user.conn != conn {
			_, err := user.conn.Write([]byte(msg + "\n"))
			if err != nil {
				log.Println("Error sending message to user:", err)
			}

		}

	}
}

func removeUser(conn net.Conn) {
	for i, user := range users {
		if user.conn == conn {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
}
