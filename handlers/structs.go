package handlers

import (
	"net"
	"sync"
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
