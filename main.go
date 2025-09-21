package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"tuff/tuff"
	"tuff/tuff/packets"
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
func handleRequest(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function exits.

	// Make a buffer to hold incoming data.
	var buf = make([]byte, 1024)

	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf[:])
	if err != nil {
		slog.Error("failed to read tcp conn")
		return
	}
	packet, err := tuff.ReadMessage(buf[:n])
	if err != nil {
		slog.Error("failed to read packet", "error", err)
		return
	}
	handshake, err := packets.DecodeHandshake(packet.Data)
	if err != nil {
		slog.Error("failed to decode handshake packet", "error", err)
		return
	}
	fmt.Printf("%+v", handshake)
}
