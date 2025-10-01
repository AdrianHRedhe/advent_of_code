package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func read_input() []string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string

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

func (pos coord) move(direction coord) coord {
	return coord{
		x: pos.x + direction.x,
		y: pos.y + direction.y,
	}
}

type state struct {
	pos   coord
	steps int
}

var (
	// y is from the top hence y: -1 is upwards
	N = coord{0, -1}
	W = coord{-1, 0}
	S = coord{0, 1}
	E = coord{1, 0}
)

func getGameMap(lines []string) (gameMap map[coord]string, start coord, end coord) {
	gameMap = make(map[coord]string)
	for y, line := range lines {
		for x, symbol := range line {
			pos := coord{x: x, y: y}
			gameMap[pos] = string(symbol)

			if string(symbol) == "S" {
				start = pos
			}
			if string(symbol) == "E" {
				end = pos
			}
		}
	}
	return gameMap, start, end
}

func nStepsShortestPath(gameMap map[coord]string, start coord, end coord) int {
	// init queue
	firstState := state{pos: start, steps: 0}
	queue := []state{firstState}
	seen := make(map[coord]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		var nextPos coord
		var nextState state

		if current.pos == end {
			return current.steps
		}

		for _, dir := range []coord{N, W, S, E} {
			nextPos = current.pos.move(dir)

			if symbol, exists := gameMap[nextPos]; exists && symbol != "#" {
				if seen[nextPos] != true {
					nextState = state{pos: nextPos, steps: current.steps + 1}
					queue = append(queue, nextState)
					seen[nextPos] = true
				}
			}
		}
	}
	return -1
}

func statesPartOfShortestPath(gameMap map[coord]string, start coord, end coord, stepsToEnd int) (states []state) {
	first := []coord{start}
	queue := [][]coord{first}
	seen := make(map[coord]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// we have gone past the limit should not be requeued
		if len(current)-1 > stepsToEnd {
			continue
		}

		currentPos := current[len(current)-1]

		if currentPos == end {
			for i, pos := range current {
				states = append(states, state{pos: pos, steps: i})
			}
		}

		for _, dir := range []coord{N, W, S, E} {
			nextPos := currentPos.move(dir)

			if symbol, exists := gameMap[nextPos]; exists && symbol != "#" {
				if seen[nextPos] != true {
					next := append(current, nextPos)
					queue = append(queue, next)
					seen[nextPos] = true
				}
			}
		}

	}
	return states
}

func calculateStepsToEndPerPos(gameMap map[coord]string, end coord) map[coord]int {
	// this is what takes time if we want to speed up this program
	stepsToEnd := make(map[coord]int)
	for pos, symbol := range gameMap {
		if symbol != "#" {
			stepsToEnd[pos] = nStepsShortestPath(gameMap, pos, end)
		}
	}
	return stepsToEnd
}

func generateManhattanPoints(center coord, radius int) []coord {
	// inspiration: https://stackoverflow.com/questions/75128474/how-to-generate-all-of-the-coordinates-that-are-within-a-manhattan-distance-r-of
	// generates points in a manhattan radius of a center
	manhattanPoints := []coord{}
	for offset := range radius {
		inverseOffset := radius - offset
		p1 := coord{x: center.x + offset, y: center.y + inverseOffset}
		p2 := coord{x: center.x + inverseOffset, y: center.y - offset}
		p3 := coord{x: center.x - offset, y: center.y - inverseOffset}
		p4 := coord{x: center.x - inverseOffset, y: center.y + offset}
		manhattanPoints = append(manhattanPoints, p1)
		manhattanPoints = append(manhattanPoints, p2)
		manhattanPoints = append(manhattanPoints, p3)
		manhattanPoints = append(manhattanPoints, p4)
	}
	return manhattanPoints
}

func teleport(gameMap map[coord]string, from coord, duration int) (to []state) {
	// jump map, from "from" you can jump to "to.pos" costing "to.steps"
	to = []state{}
	for radius := 1; radius <= duration; radius++ {
		for _, point := range generateManhattanPoints(from, radius) {
			if symbol, exists := gameMap[point]; exists && symbol != "#" {
				after := state{pos: point, steps: radius}
				to = append(to, after)
			}
		}
	}
	return to
}

func countCheatsThatSave100Steps(gameMap map[coord]string, start coord, end coord, cheatDuration int, stepsToEnd map[coord]int) (nCheats int) {
	nCheats = 0
	stepThreshold := 100
	numberOfStepsShortestPath := nStepsShortestPath(gameMap, start, end)
	potentialCheatStates := statesPartOfShortestPath(gameMap, start, end, numberOfStepsShortestPath)

	for _, from := range potentialCheatStates {
		for _, to := range teleport(gameMap, from.pos, cheatDuration) {
			cheatSteps := from.steps + to.steps + stepsToEnd[to.pos]
			if (numberOfStepsShortestPath - cheatSteps) >= stepThreshold {
				nCheats++
			}
		}
	}

	return nCheats
}

func main() {
	lines := read_input()
	gameMap, start, end := getGameMap(lines)
	stepsToEnd := calculateStepsToEndPerPos(gameMap, end)

	cheatDuration := 2
	nCheats := countCheatsThatSave100Steps(gameMap, start, end, cheatDuration, stepsToEnd)
	fmt.Println("part 1: ", nCheats)

	cheatDuration = 20
	nCheats = countCheatsThatSave100Steps(gameMap, start, end, cheatDuration, stepsToEnd)
	fmt.Println("part 2: ", nCheats)
}
