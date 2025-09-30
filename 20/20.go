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

func getGameMap(lines []string) (gameMap map[coord]string, start coord, end coord, boundary coord) {
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
	boundary = coord{x: len(lines[0]), y: len(lines)}
	return gameMap, start, end, boundary
}

var (
	// y is from the top hence y: -1 is upwards
	N = coord{0, -1}
	W = coord{-1, 0}
	S = coord{0, 1}
	E = coord{1, 0}
)

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

func findPotentialCheats(gameMap map[coord]string, boundary coord) []coord {
	xBound, yBound := boundary.x, boundary.y
	potentialCheats := []coord{}

	for x := range xBound {
		for y := range yBound {
			// if x is on the edges then we dont want to cheat
			if x == 0 || x == xBound {
				continue
			}
			// if y is on the edges then we dont want to cheat
			if y == 0 || y == yBound {
				continue
			}
			cheat := coord{x: x, y: y}
			// if it is not a # then it is not a cheat
			if symbol, exists := gameMap[cheat]; !exists || symbol != "#" {
				continue
			}
			potentialCheats = append(potentialCheats, cheat)
		}
	}
	return potentialCheats
}

func countCheatsThatSave100Steps(gameMap map[coord]string, start coord, end coord, boundary coord) (nCheats int) {
	var cheatGameMap map[coord]string
	var cheatSteps int
	nCheats = 0
	nStepsOriginal := nStepsShortestPath(gameMap, start, end)
	stepThreshold := 100

	potentialCheats := findPotentialCheats(gameMap, boundary)
	for _, cheat := range potentialCheats {
		cheatGameMap = make(map[coord]string)
		for k, v := range gameMap {
			cheatGameMap[k] = v
		}
		cheatGameMap[cheat] = "."
		cheatSteps = nStepsShortestPath(cheatGameMap, start, end)
		if (nStepsOriginal - cheatSteps) >= stepThreshold {
			nCheats++
		}
	}

	return nCheats
}

func main() {
	lines := read_input()
	gameMap, start, end, boundary := getGameMap(lines)
	nCheats := countCheatsThatSave100Steps(gameMap, start, end, boundary)

	fmt.Println("part 1: ", nCheats)

	// fmt.Println("part 2: ", )
}
