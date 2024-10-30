// server.go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type World struct {
	Name    string
	Players map[int]*Player
	mu      sync.Mutex
}

type Player struct {
	ID         int
	Name       string
	Level      int
	Experience int
	Coordinate Coordinate
}

type Coordinate struct {
	X float64
	Y float64
}

type MoveCommand struct {
	PlayerID int
	Direction string // "up", "down", "left", "right"
}

var world = World{
	Name:    "Fantasy Land",
	Players: make(map[int]*Player),
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("MMORPG Server started on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var player Player
	if err := json.NewDecoder(conn).Decode(&player); err != nil {
		fmt.Println("Error decoding player data:", err)
		return
	}

	world.mu.Lock()
	world.Players[player.ID] = &player
	world.mu.Unlock()

	fmt.Printf("Player %s entered the world at coordinates: %+v\n", player.Name, player.Coordinate)

	for {
		var command MoveCommand
		if err := json.NewDecoder(conn).Decode(&command); err != nil {
			fmt.Println("Error decoding move command:", err)
			return
		}

		world.mu.Lock()
		player := world.Players[command.PlayerID]
		switch command.Direction {
		case "up":
			player.Coordinate.Y += 1.0
		case "down":
			player.Coordinate.Y -= 1.0
		case "left":
			player.Coordinate.X -= 1.0
		case "right":
			player.Coordinate.X += 1.0
		}
		world.mu.Unlock()

		fmt.Printf("Player %s moved %s to %+v\n", player.Name, command.Direction, player.Coordinate)
	}
}
