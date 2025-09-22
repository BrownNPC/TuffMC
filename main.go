package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"
	"tuff/connection"
	"tuff/packet"
)

func main() {
	l, err := net.Listen("tcp", "localhost:25565")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// handleRequest handles incoming requests from clients.
func handleRequest(socket net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from error", "error", r)
		}
	}()

	defer socket.Close() // Ensure the connection is closed when the function exits.

	conn := connection.NewConnection(socket)
	_, err := conn.ReadMsg(1)
	if err != nil {
		fmt.Println(err)
	}
}
