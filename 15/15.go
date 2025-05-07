package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func read_input() (lines []string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

type coord struct {
	x, y int
}

func (this *coord) move(movement coord) coord {
	return coord{
		x: this.x + movement.x,
		y: this.y + movement.y,
	}
}

func formatInput(lines []string) (warehouse map[coord]string, movements []string) {
	var idxEmpty int
	warehouse = make(map[coord]string)
	// first part of input will contain information on
	// warehouse and end with an empty line
	for i := range lines {
		if lines[i] == "" {
			idxEmpty = i
			break
		}
		for j := range lines[i] {
			position := coord{i, j}
			warehouse[position] = string(lines[i][j])
		}
	}
	// second part of input will contain information on
	// movements
	for _, line := range lines[idxEmpty+1:] {
		for _, char := range line {
			movements = append(movements, string(char))
		}
	}
	return warehouse, movements
}

func canMove(warehouse map[coord]string, position coord, movement coord, nMoves int) (bool, int) {
	// returns true if next position is empty recursively
	// returns false if next position is cannot be moved
	// recursively
	nMoves++
	nextPosition := position.move(movement)
	nextObject := warehouse[nextPosition]
	if nextObject == "." {
		return true, nMoves
	}
	if nextObject == "#" {
		return false, 0
	}
	if nextObject == "O" {
		return canMove(warehouse, nextPosition, movement, nMoves)
	}
	log.Fatal("Next object is neither . # or O")
	return false, 0
}

func update(warehouse map[coord]string, position coord, movement string) (updatedWarehouse map[coord]string, guardPosition coord) {
	// Identify the movement
	movementTypes := map[string]coord{
		"^": {-1, 0},
		">": {0, 1},
		"v": {1, 0},
		"<": {0, -1},
	}
	movementCoord := movementTypes[movement]
	if makeMove, nMoves := canMove(warehouse, position, movementCoord, 0); makeMove {
		curPos := position
		nextPos := position.move(movementCoord)

		// Will be empty were guard was
		warehouse[curPos] = "."
		warehouse[nextPos] = "@"
		guardPosition = nextPos
		for range nMoves - 1 {
			curPos = nextPos
			nextPos = curPos.move(movementCoord)
			// We can only move along Os
			warehouse[nextPos] = "O"
		}
	} else {
		guardPosition = position
	}
	return warehouse, guardPosition
}

func calculateScore(warehouse map[coord]string) (score int) {
	for position, object := range warehouse {
		if object != "O" {
			continue
		}
		score += (position.x*100 + position.y)
	}
	return score
}

func printWarehouse(warehouse map[coord]string) {
	// Only used for printing first example
	println()
	h, w := 8, 9
	matrix := make([][]string, h)
	for i := range h {
		matrix[i] = make([]string, w)
	}
	for k, v := range warehouse {
		matrix[k.x][k.y] = string(v)
	}
	for i := range h {
		for j := range w {
			fmt.Print(string(matrix[i][j]))
		}
		println()
	}
}

func main() {
	lines := read_input()
	warehouse, movements := formatInput(lines)

	// part 1:
	// find guard
	var guardPosition coord
	for k, v := range warehouse {
		if v == "@" {
			guardPosition = k
		}
	}

	// update warehouse for each movement
	for _, movement := range movements {
		// print(movement)
		warehouse, guardPosition = update(warehouse, guardPosition, movement)
		// printWarehouse(warehouse)
	}

	// calculate score
	score := calculateScore(warehouse)
	fmt.Println("part 1: ", score)
}
