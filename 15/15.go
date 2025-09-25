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

func formatInputPart2(lines []string) (warehouse map[coord]string, movements []string) {
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
			position_1 := coord{i, 2 * j}
			position_2 := coord{i, 2*j + 1}
			switch string(lines[i][j]) {
			case "#":
				warehouse[position_1] = "#"
				warehouse[position_2] = "#"
			case ".":
				warehouse[position_1] = "."
				warehouse[position_2] = "."
			case "O":
				warehouse[position_1] = "["
				warehouse[position_2] = "]"
			case "@":
				warehouse[position_1] = "@"
				warehouse[position_2] = "."
			default:
				log.Fatal("Warehouse object not in #.O@")
			}
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

func canMoveAndCollectPart2(warehouse map[coord]string, position coord, movement coord, positionsToMove map[coord]bool) bool {
	if positionsToMove[position] {
		return true // Already processed this position
	}

	nextPos := position.move(movement)
	nextCell := warehouse[nextPos]

	switch nextCell {
	case "#":
		return false
	case ".":
		positionsToMove[position] = true
		return true
	case "[":
		// Vertical movement - Moving into left part of box
		rightPos := coord{nextPos.x, nextPos.y + 1}
		if movement.y == 0 { // Vertical movement - need to move both parts
			canMoveLeft := canMoveAndCollectPart2(warehouse, nextPos, movement, positionsToMove)
			canMoveRight := canMoveAndCollectPart2(warehouse, rightPos, movement, positionsToMove)
			if canMoveLeft && canMoveRight {
				positionsToMove[position] = true
				return true
			}
			return false
		} else { // Horizontal movement
			if canMoveAndCollectPart2(warehouse, nextPos, movement, positionsToMove) {
				positionsToMove[position] = true
				return true
			}
			return false
		}
	case "]":
		// Vertical movement - Moving into right part of box
		leftPos := coord{nextPos.x, nextPos.y - 1}
		if movement.y == 0 { // Vertical movement - need to move both parts
			canMoveLeft := canMoveAndCollectPart2(warehouse, leftPos, movement, positionsToMove)
			canMoveRight := canMoveAndCollectPart2(warehouse, nextPos, movement, positionsToMove)
			if canMoveLeft && canMoveRight {
				positionsToMove[position] = true
				return true
			}
			return false
		} else { // Horizontal movement
			if canMoveAndCollectPart2(warehouse, nextPos, movement, positionsToMove) {
				positionsToMove[position] = true
				return true
			}
			return false
		}
	default:
		log.Fatal("Next object is neither . # [ or ]")
		return false
	}
}

func updatePart2(warehouse map[coord]string, position coord, movement string) (map[coord]string, coord) {
	movementTypes := map[string]coord{
		"^": {-1, 0},
		">": {0, 1},
		"v": {1, 0},
		"<": {0, -1},
	}
	movementCoord := movementTypes[movement]

	positionsToMove := make(map[coord]bool)
	if canMoveAndCollectPart2(warehouse, position, movementCoord, positionsToMove) {
		// Create new warehouse state
		newWarehouse := make(map[coord]string)
		for k, v := range warehouse {
			newWarehouse[k] = v
		}

		// Clear old positions
		// positionsToMove is actually a pointer and will be updated in canMoveAndCollectPart2
		// TODO: make the return return of this value more explicit.
		for pos := range positionsToMove {
			newWarehouse[pos] = "."
		}

		// Move everything to new positions
		for pos := range positionsToMove {
			newPos := pos.move(movementCoord)
			newWarehouse[newPos] = warehouse[pos]
		}

		return newWarehouse, position.move(movementCoord)
	}

	return warehouse, position
}

func calculateScorePart2(warehouse map[coord]string) (score int) {
	for position, object := range warehouse {
		if object != "[" {
			continue
		}
		score += (position.x*100 + position.y)
	}
	return score
}

func main() {
	lines := read_input()

	// part 1:
	// init
	warehouse, movements := formatInput(lines)

	// find guard
	var guardPosition coord

	for k, v := range warehouse {
		if v == "@" {
			guardPosition = k
		}
	}

	// update warehouse for each movement
	for _, movement := range movements {
		warehouse, guardPosition = update(warehouse, guardPosition, movement)
	}

	// calculate score
	score := calculateScore(warehouse)
	fmt.Println("part 1: ", score)

	// part 2:
	// Init
	warehouse, movements = formatInputPart2(lines)

	// find guard
	for k, v := range warehouse {
		if v == "@" {
			guardPosition = k
		}
	}

	// // update warehouse for each movement
	for _, movement := range movements {
		warehouse, guardPosition = updatePart2(warehouse, guardPosition, movement)
	}

	// calculate score
	score = calculateScorePart2(warehouse)
	fmt.Println("part 2: ", score)
}
