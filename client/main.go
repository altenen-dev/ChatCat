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

func main() {
	port := "8989" // Default port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal("Error resolving TCP address: ", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		log.Fatal("Error dialing TCP: ", err)
	}
	defer conn.Close()

	// Prompt the user for their name
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input: ", err)
	}
	clientName = strings.TrimSpace(clientName)

	// Send the user's name to the server
	_, err = conn.Write([]byte(clientName + "\n"))
	if err != nil {
		log.Fatal("Error writing the name: ", err)
	}

	go receiveMessages(conn)

	// Start the chat loop
	for {
		fmt.Print("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + clientName + "]: ")
		clientMsg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input: ", err)
		}
		fmt.Println()
		clientMsg = strings.TrimSpace(clientMsg)

		if clientMsg != "" {
			_, err = conn.Write([]byte(clientMsg + "\n"))
			if err != nil {
				log.Fatal("Error writing the message: ", err)
			}
		}
	}
}

func receiveMessages(conn *net.TCPConn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from server:", err)
			return
		}
		fmt.Println(string(buffer[:n]))
	}
}
