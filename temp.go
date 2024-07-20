package main

import (
	"fmt"
	"log"
	"net"
	handle "net-cat/handlers"
)

func main() {
	port := "8989" // Default port

	dataStream, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TCP server is running on:", port)
	defer dataStream.Close()

	go handle.BroadcastMessages() // Start a goroutine to broadcast messages
	for {
		conn, err := dataStream.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		if len(handle.Clients) >= 10 {
			conn.Write([]byte("Max connections reached, rejecting new connection\n"))
			conn.Close()
			continue
		}

		client, err := handle.NewClientsHandler(conn)
		if err != nil {
			log.Println("Error reading client name:", err)
			conn.Close()
			continue
		}

		go handle.HandleClient(client)
	}
}
