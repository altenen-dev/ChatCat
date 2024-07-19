package main

import "net"

type User struct {
	name string
	conn net.Conn
}
