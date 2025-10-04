package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"tuff/connection"
	"tuff/packet"

	"github.com/coder/websocket"
)

func main() {
	go startWebsocketListener("localhost:8081")
	startTCPListener("localhost:25565")
}

// EaglerCraft v1 handshake protocol packet
type EaglerHandshakePacket struct {
	// Must be 1
	PacketId byte
	// Must be 1
	EaglerCraftVersion byte
	MinecraftVersion   byte

	// 1 byte prefixed for string length.
	ClientBrand string
	// 1 byte prefixed for string length
	ClientVersion string
}

// EaglerCraft v1 handshake protocol response.
type EaglerHandshakeResponse struct {
	// Must be 2
	PacketID byte
	// Must be 1
	EaglerCraftVersion byte
	// 1 byte prefixed for string length.
	ServerName string
	// 1 byte prefixed for string length.
	ServerVersion string

	// 3 bytes of padding required for some reason
	// https://github.com/ayunami2000/ayunViaProxyEagUtils/blob/main/src/main/java/me/ayunami2000/ayunViaProxyEagUtils/EaglercraftHandler.java#L218
	_ [3]byte
}

// handleRequest handles incoming requests from clients.
func handleRequest(conn *connection.Connection) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from error", "error", r)
			fmt.Println(string(debug.Stack()))
		}
	}()
	defer conn.Close()

	err := conn.HandleHandshake(packet.StatusResponsePacketConfig{
		PlayerCount: 0,
		Description: "Tuff server gng",
	})
	if err != nil {
		slog.Error("could not handle handshake", "error", err)
		return
	}
	// This could be because the client did not want to login (status ping)
	if !conn.IsLoggedIn() {
		return
	}

	err = conn.WriteMessage(packet.EncodeJoinGamePacket(packet.JoinGamePacketConfig{
		EntityID:         0,
		Gamemode:         0,
		Dimension:        0,
		Difficulty:       1,
		MaxPlayers:       0,
		LevelType:        "default",
		ReducedDebugInfo: false,
	}))
	if err != nil {
		slog.Error("failed to send join game packet", "error", err)
		return
	}
	err = conn.WriteMessage(packet.EncodeSpawnPositionPacket(0, 100, 0))
	if err != nil {
		slog.Error("failed to send spawn position packet", "error", err)
		return
	}
	err = conn.WriteMessage(packet.EncodePlayerPositionAndLookPacket(packet.PlayerPositionAndLookConfig{
		X:          0,
		Y:          0,
		Z:          0,
		Pitch:      0,
		Yaw:        0,
		Flags:      0,
		TeleportId: 0,
	}))
	if err != nil {
		slog.Error("failed to send Player Position and Look packet", "error", err)
		return
	}
}
func startTCPListener(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "error", err)
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(connection.NewConnection(conn))
	}
}

// Maybe support eagler craft in the future?
func startWebsocketListener(addr string) {
	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		ws.Write(context.Background(), websocket.MessageText, []byte(jzml))
		fmt.Println("aaa")
		if err != nil {
			slog.Error("failed to accept websocket connection", "error", err)
			return
		}
		conn := websocket.NetConn(context.Background(), ws, websocket.MessageBinary)
		go handleRequest(connection.NewConnection(conn))
	}))
}

const jzml = `{"name":"EaglercraftXServer","brand":"lax1dude","vers":"EaglercraftXServer/1.0.7","plaf":"velocity","cracked":true,"time":1758554003989,"uuid":"cb703b97-1a8f-4a77-b6f9-ba16e9d8a4c4","type":"motd","data":{"cache":true,"motd":["§a§lWebMC","§7Minecraft Oneblock... With a twist!§r"],"icon":true,"online":32,"max":250,"players":[]}}`
