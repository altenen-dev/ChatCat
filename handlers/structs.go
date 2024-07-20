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
	Clients = make(map[net.Conn]Client)
	Msgs    = make(chan Message) // Initialize the channel
	Mutex   sync.Mutex
	History []string
)
