package main

import (
	"bufio"
	"container/heap"
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

var (
	N = coord{0, 1}
	E = coord{1, 0}
	S = coord{0, -1}
	W = coord{-1, 0}
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

type reindeer struct {
	pos       coord
	direction coord
}

func (current reindeer) move() reindeer {
	return reindeer{
		current.pos.move(current.direction),
		current.direction,
	}
}

func (current reindeer) turnClockwise() reindeer {
	var newDirection coord
	switch current.direction {
	case E:
		newDirection = S
	case S:
		newDirection = W
	case W:
		newDirection = N
	case N:
		newDirection = E
	}

	if (newDirection != N) && (newDirection != E) && (newDirection != S) && (newDirection != W) {
		log.Fatal("current direction not in NESW")
	}

	return reindeer{
		current.pos,
		newDirection,
	}
}

func (current reindeer) turnCounterClockwise() reindeer {
	var newDirection coord
	switch current.direction {
	case E:
		newDirection = N
	case S:
		newDirection = E
	case W:
		newDirection = S
	case N:
		newDirection = W
	}

	if (newDirection != N) && (newDirection != E) && (newDirection != S) && (newDirection != W) {
		log.Fatal("current direction not in NESW")
	}

	return reindeer{
		current.pos,
		newDirection,
	}
}

// heap implementation using an interface
type State struct {
	reindeer reindeer
	cost     int
	path     []coord
}
type PriorityQueue []*State

// required by sort interface (in turn required by heap interface)
func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// required by heap interface
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*State))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	item := old[len(old)-1]   // topmost item
	*pq = old[0 : len(old)-1] // heap except topmost item
	return item
}

func createGameMap(lines []string) (map[coord]string, coord, coord) {
	gameMap := make(map[coord]string)

	var start, end coord

	// since we are reading from upside down we should go
	// from highest to lowest y value
	y_max := len(lines)

	for i := range lines {
		line := lines[i]
		y := y_max - i

		for j, symbol := range line {
			x := j
			pos := coord{x, y}
			gameMap[pos] = string(symbol)

			if string(symbol) == "S" {
				start = pos
			} else if string(symbol) == "E" {
				end = pos
			}

		}
	}

	return gameMap, start, end
}

func calculateCostOfOptimalPath(gameMap map[coord]string, start coord, end coord) (minimumCost int) {
	// initialize the reindeer
	startReindeer := reindeer{pos: start, direction: E}
	startState := State{reindeer: startReindeer, cost: 0, path: nil}

	// init priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &startState)

	// save minimum cost per location
	visisted := make(map[reindeer]int)

	for pq.Len() > 0 {
		var nextReindeer reindeer
		var newCost int

		current := heap.Pop(pq).(*State)

		// if we have reached the end return the cost
		if current.reindeer.pos == end {
			return current.cost
		}

		// skip if we have have already seen this pos / direction at a lower cost
		if prevCost, exists := visisted[current.reindeer]; exists && prevCost < current.cost {
			continue
		}
		// update cost since this reindeer is either seen for the first time or lower cost than last time.
		visisted[current.reindeer] = current.cost

		// try all actions for this state
		// move forward
		nextReindeer = current.reindeer.move()
		// can only move forward if there is no wall
		if gameMap[nextReindeer.pos] != "#" {
			newCost = current.cost + 1
			// push to heap if reindeer is seen for first time or lower cost than last time.
			if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {
				heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost})
			}
		}

		// turn clockwise
		nextReindeer = current.reindeer.turnClockwise()
		newCost = current.cost + 1000
		// push to heap if reindeer is seen for first time or lower cost than last time.
		if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {
			heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost})
		}

		// turn counter clockwise
		nextReindeer = current.reindeer.turnCounterClockwise()
		newCost = current.cost + 1000
		// push to heap if reindeer is seen for first time or lower cost than last time.
		if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {
			heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost})
		}

	}

	return -1
}

func numberOfTilesPartOfOptimalPath(gameMap map[coord]string, start coord, end coord) (n_tiles int) {
	optimalPathCost := calculateCostOfOptimalPath(gameMap, start, end)
	// if there is no path to win, then 0 tiles were part of optimal path
	if optimalPathCost == -1 {
		n_tiles = 0
		return n_tiles
	}

	// initialize the reindeer
	startReindeer := reindeer{pos: start, direction: E}
	startState := State{reindeer: startReindeer, cost: 0, path: []coord{start}}

	// init priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &startState)

	// save minimum cost per location
	visisted := make(map[reindeer]int)

	// save tiles part of optimal paths
	optimalTiles := make(map[coord]bool)

	for pq.Len() > 0 {
		var nextReindeer reindeer
		var newCost int
		var newPath []coord

		current := heap.Pop(pq).(*State)

		// If cost is higher than optimal path then we can skip it.
		if current.cost > optimalPathCost {
			break
		}

		// If we have reached the end and cost is optimal these tiles are part of an optimal path.
		if (current.reindeer.pos == end) && (current.cost == optimalPathCost) {
			for _, tile := range current.path {
				optimalTiles[tile] = true
			}
			continue
		}
		// We have visited this State at a cheaper price => suboptimal => skip it.
		if prevCost, exists := visisted[current.reindeer]; exists && prevCost < current.cost {
			continue
		}
		// update lowest known cost for this location since this reindeer is either seen for the
		// first time or lower cost than last time.
		visisted[current.reindeer] = current.cost

		// try all actions for this state
		// move forward
		nextReindeer = current.reindeer.move()
		// can only move forward if there is no wall
		if gameMap[nextReindeer.pos] != "#" {
			newCost = current.cost + 1
			// push to heap if reindeer is seen for first time or lower cost than last time.
			if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {

				// also update paths which have been visited.
				// Pre allocate another space in order to avoid having to move it over when we append.
				newPath = make([]coord, len(current.path), len(current.path)+1)
				// copy over current values
				copy(newPath, current.path)
				newPath = append(newPath, nextReindeer.pos)

				heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost, path: newPath})
			}
		}
		// turn clockwise
		nextReindeer = current.reindeer.turnClockwise()
		newCost = current.cost + 1000
		// push to heap if reindeer is seen for first time or lower cost than last time.
		if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {
			heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost, path: current.path})
		}
		// turn counter clockwise
		nextReindeer = current.reindeer.turnCounterClockwise()
		newCost = current.cost + 1000
		// push to heap if reindeer is seen for first time or lower cost than last time.
		if prevCost, exists := visisted[nextReindeer]; !exists || current.cost < prevCost {
			heap.Push(pq, &State{reindeer: nextReindeer, cost: newCost, path: current.path})
		}
	}

	return len(optimalTiles)
}

func main() {
	lines := read_input()
	gameMap, start, end := createGameMap(lines)
	minScore := calculateCostOfOptimalPath(gameMap, start, end)

	fmt.Println("part 1: ", minScore)

	n_tiles := numberOfTilesPartOfOptimalPath(gameMap, start, end)
	fmt.Println("part 2: ", n_tiles)
}
