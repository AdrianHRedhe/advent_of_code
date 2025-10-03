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

func (pos coord) move(dir coord) coord {
	return coord{
		x: pos.x + dir.x,
		y: pos.y + dir.y,
	}
}

type state struct {
	from  coord
	to    coord
	depth int
}

var (
	up    = coord{0, 1}
	down  = coord{0, -1}
	left  = coord{-1, 0}
	right = coord{1, 0}
)

var (
	symbolToDirection = map[string]coord{
		"^": up,
		"v": down,
		"<": left,
		">": right,
	}
	directionToSymbol = map[coord]string{
		up:    "^",
		down:  "v",
		left:  "<",
		right: ">",
	}
)

var (
	numericKeypad = map[coord]string{
		{0, 3}: "7", {1, 3}: "8", {2, 3}: "9",
		{0, 2}: "4", {1, 2}: "5", {2, 2}: "6",
		{0, 1}: "1", {1, 1}: "2", {2, 1}: "3",
		{1, 0}: "0", {2, 0}: "A",
	}
	directionalKeypad = map[coord]string{
		{1, 1}: "^", {2, 1}: "A",
		{0, 0}: "<", {1, 0}: "v", {2, 0}: ">",
	}
)

var memo = make(map[state]int)

func read_input() ([]string, []int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var codes []string
	var values []int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		codes = append(codes, code)

		numeric := code[:len(code)-1]

		value, err := strconv.Atoi(numeric)
		if err != nil {
			log.Println("could not convert code to int", err)
			continue
		}
		values = append(values, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return codes, values
}

func isValidStepSequence(start coord, sequence string, keyPad map[coord]string) bool {
	current := start
	for _, symbol := range sequence {
		movement := symbolToDirection[string(symbol)]
		next := current.move(movement)
		if _, exists := keyPad[next]; !exists {
			return false
		}
		current = next
	}
	return true
}

func shortestPath(start coord, end coord, keyPad map[coord]string) (subsolutions []string) {
	// base case
	if start == end {
		return []string{"A"}
	}

	xDelta := end.x - start.x
	yDelta := end.y - start.y

	var horizontalMovement string
	var verticalMovement string
	if xDelta < 0 {
		horizontalMovement = strings.Repeat("<", -xDelta)
	} else {
		horizontalMovement = strings.Repeat(">", xDelta)
	}
	if yDelta < 0 {
		verticalMovement = strings.Repeat("v", -yDelta)
	} else {
		verticalMovement = strings.Repeat("^", yDelta)
	}

	solution1 := horizontalMovement + verticalMovement + "A"
	solution2 := verticalMovement + horizontalMovement + "A"

	if isValidStepSequence(start, solution1, keyPad) {
		subsolutions = append(subsolutions, solution1)
	}
	if solution1 != solution2 && isValidStepSequence(start, solution2, keyPad) {
		subsolutions = append(subsolutions, solution2)
	}
	return subsolutions
}

func shortestPathLength(start coord, end coord, keypad map[coord]string, depth int) int {
	key := state{start, end, depth}
	// return using memo if youve seen this
	if minLen, exists := memo[key]; exists {
		return minLen
	}

	sequences := shortestPath(start, end, keypad)

	if depth == 0 {
		minLen := 0
		// base case return len short path
		if len(sequences) == 1 || len(sequences[0]) < len(sequences[1]) {
			minLen = len(sequences[0])
		} else {
			minLen = len(sequences[1])
		}
		memo[key] = minLen
		return minLen
	}

	minCost := int(^uint(0) >> 1) // Max int
	for _, sequence := range sequences {
		cost := calculateSequenceCost(sequence, directionalKeypad, depth-1)
		if cost < minCost {
			minCost = cost
		}
	}

	memo[key] = minCost
	return minCost
}

func calculateSequenceCost(sequence string, keypad map[coord]string, depth int) (totalCost int) {
	positions := map[string]coord{}
	for k, v := range keypad {
		positions[v] = k
	}

	currentPos := positions["A"]
	for _, symbol := range sequence {
		targetPos := positions[string(symbol)]
		totalCost += shortestPathLength(currentPos, targetPos, keypad, depth)
		currentPos = targetPos
	}

	return totalCost
}

func main() {
	codes, values := read_input()
	complexity := 0
	for i, code := range codes {
		numericCode := values[i]
		cost := calculateSequenceCost(code, numericKeypad, 2)
		complexity += cost * numericCode
	}

	fmt.Println("part 1: ", complexity)

	complexity = 0
	for i, code := range codes {
		numericCode := values[i]
		cost := calculateSequenceCost(code, numericKeypad, 25)
		complexity += cost * numericCode
	}
	fmt.Println("part 2: ", complexity)
}
