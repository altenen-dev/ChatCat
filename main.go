package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	handle "net-cat/handlers"
	"os"
	"os/signal"
	"syscall"
)

const defaultPort = "8989"

func main() {
	port := flag.String("port", defaultPort, "Port to run the TCP server on")
	flag.Parse()

	dataStream, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TCP server is running on port:", *port)
	defer dataStream.Close()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nShutting down server...")
		dataStream.Close()
		os.Exit(0)
	}()

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

		go func(conn net.Conn) {
			client, err := handle.NewClientsHandler(conn)
			if err != nil {
				log.Println("Error reading client name:", err)
				conn.Close()
				return
			}
			handle.HandleClient(client)
		}(conn)
	}
}
