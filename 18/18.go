package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

func (current coord) move(movement coord) coord {
	return coord{
		current.x + movement.x,
		current.y + movement.y,
	}
}

var (
	// North is -1 since y coords are going in opposite direction
	N = coord{0, -1}
	E = coord{1, 0}
	S = coord{0, 1}
	W = coord{-1, 0}
)

func read_input() []coord {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var corrupted []coord

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 2 {
			log.Printf("wrong input for line: %s", line)
			continue
		}

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("could not convert %s to %v", parts[0], err)
			continue
		}

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("could not convert %s to %v", parts[1], err)
			continue
		}
		corruptedCord := coord{x: x, y: y}
		corrupted = append(corrupted, corruptedCord)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return corrupted
}

func createMemorySpace() (memorySpace map[coord]string, start coord, end coord) {
	size := 71
	memorySpace = make(map[coord]string)
	start = coord{0, 0}
	end = coord{size - 1, size - 1}
	for i := range size {
		for j := range size {
			memorySpace[coord{i, j}] = "."
		}
	}
	return memorySpace, start, end
}

func corruptMemorySpace(corrupted []coord, memorySpace map[coord]string) map[coord]string {
	kilobyte := 1024
	for i := range kilobyte {
		nextCorruptedByte := corrupted[i]
		memorySpace[nextCorruptedByte] = "#"
	}
	return memorySpace
}

// Struct to store the state
type State struct {
	pos   coord
	steps int
}

func calculateShortestPath(corruptedMemorySpace map[coord]string, start coord, end coord) int {
	// initialize the state
	startState := State{pos: start, steps: 0}
	queue := []State{startState}

	// save seen locations
	visited := make(map[coord]bool)
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// we have reached the end and should therefore return the number of steps
		if current.pos == end {
			return current.steps
		}

		// try all possible movements
		for _, direction := range []coord{N, W, S, E} {
			nextPos := current.pos.move(direction)

			if cellType, exists := corruptedMemorySpace[nextPos]; exists && cellType != "#" {
				if !visited[nextPos] {
					visited[nextPos] = true
					// push new to queue
					nextState := State{pos: nextPos, steps: current.steps + 1}
					queue = append(queue, nextState)
				}
			}
		}

	}

	return -1
}

func main() {
	corrupted := read_input()
	memorySpace, start, end := createMemorySpace()
	corruptedMemorySpace := corruptMemorySpace(corrupted, memorySpace)
	shortestPath := calculateShortestPath(corruptedMemorySpace, start, end)

	fmt.Println("part 1: ", shortestPath)
}
