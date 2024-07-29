package main

import (
	"fmt"
	"log"
	"net"
	handle "net-cat/handlers"
	"os"
)

const defaultPort = "8989"

func main() {
	port := defaultPort
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	var localhostIP string
	if len(os.Args) == 2 {
		port = os.Args[1]
		localhostIP := net.IPv4(127, 0, 0, 1).String()
		fmt.Println("Localhost IP address:", localhostIP)
	}

	//Starting tcp connection on $port
	listener, err := net.Listen("tcp", localhostIP+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TCP server is running on port:", port)
	defer listener.Close()

	// Start a goroutine to broadcast messages
	go handle.BroadcastMessages()

	//Start accepting clients connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			conn.Close()
			continue
		}

		if len(handle.Clients) >= 10 {
			conn.Write([]byte("Max connections reached, please try again later\n"))
			conn.Close()
			continue
		}
		//Start go routine to handle clients and start chatting
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
