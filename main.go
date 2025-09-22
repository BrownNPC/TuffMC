package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
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
func handleRequest(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from error", "error", r)
		}
	}()

	defer conn.Close() // Ensure the connection is closed when the function exits.

	// Make a buffer to hold incoming data.
	var buf = make([]byte, 1024)

	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf[:])
	if err != nil {
		slog.Error("failed to read tcp conn")
		return
	}
	msg, err := packet.DecodeMessage(buf[:n])
	if err != nil {
		slog.Error("failed to read packet", "error", err)
		return
	}
	handshake, err := packet.DecodeHandshake(msg.Data)
	if err != nil {
		slog.Error("failed to decode handshake", "error", err)
		return
	}
	fmt.Printf("%+v", handshake)

	msg = packet.EncodeStatusResponse(2, "Sigma", "")
	encodedMsg := packet.EncodeMessage(msg)
	_, err = conn.Write(encodedMsg)
	if err != nil {
		slog.Error("failed to write to socket", "error", err)
		return
	}
}
