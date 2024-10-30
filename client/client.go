// client.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

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
	PlayerID  int
	Direction string // "up", "down", "left", "right"
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	player := Player{
		ID:    1,
		Name:  "Player1",
		Level: 1,
		Coordinate: Coordinate{
			X: 0.0, Y: 0.0,
		},
	}

	if err := json.NewEncoder(conn).Encode(player); err != nil {
		fmt.Println("Error sending player data:", err)
		return
	}

	fmt.Println("Connected to server as:", player.Name)
	fmt.Println("Use WASD to move (W: up, A: left, S: down, D: right)")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		direction := strings.TrimSpace(input)

		var moveCommand MoveCommand
		moveCommand.PlayerID = player.ID

		switch direction {
		case "w":
			moveCommand.Direction = "up"
		case "s":
			moveCommand.Direction = "down"
		case "a":
			moveCommand.Direction = "left"
		case "d":
			moveCommand.Direction = "right"
		default:
			fmt.Println("Invalid command. Use W, A, S, D to move.")
			continue
		}

		if err := json.NewEncoder(conn).Encode(moveCommand); err != nil {
			fmt.Println("Error sending move command:", err)
			return
		}
	}
}
