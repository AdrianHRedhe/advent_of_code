package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Coord struct {
	x, y int
}

func readInput(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func getLabMap(lines []string) map[Coord]string {
	labMap := make(map[Coord]string)
	for x, line := range lines {
		chars := strings.Split(line, "")
		for y, char := range chars {
			labMap[Coord{x, y}] = char
		}
	}
	return labMap
}

func moveGuard(labMap map[Coord]string) (map[Coord]string, bool) {
	guardStates := []string{"^", ">", "v", "<"}
	movementStates := []Coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	for coord, value := range labMap {
		if slices.Contains(guardStates, value) {
			guardIndex := slices.Index(guardStates, value)
			movement := movementStates[guardIndex]
			nextCoord := Coord{coord.x + movement.x, coord.y + movement.y}

			if _, isInside := labMap[nextCoord]; !isInside {
				labMap[coord] = "X"
				return labMap, true
			}

			if labMap[nextCoord] == "#" {
				nextGuardIndex := (guardIndex + 1) % len(guardStates)
				labMap[coord] = guardStates[nextGuardIndex]
				return labMap, false
			}

			labMap[nextCoord] = labMap[coord]
			labMap[coord] = "X"
			return labMap, false
		}
	}
	return labMap, false
}

func countVisitedLocations(labMap map[Coord]string) int {
	count := 0
	for _, value := range labMap {
		if value == "X" {
			count++
		}
	}
	return count
}

func getGuardPosition(labMap map[Coord]string) (Coord, string) {
	guardStates := []string{"^", ">", "v", "<"}
	for coord, value := range labMap {
		if slices.Contains(guardStates, value) {
			return coord, value
		}
	}
	return Coord{-1, -1}, ""
}

func hasLoop(labMap map[Coord]string, obstacleCoord Coord) bool {
	tempMap := make(map[Coord]string)
	for k, v := range labMap {
		tempMap[k] = v
	}
	tempMap[obstacleCoord] = "#"

	previousGuardStates := make(map[Coord][]string)
	isCompleted := false

	for !isCompleted {
		guardCoord, guardState := getGuardPosition(tempMap)

		if slices.Contains(previousGuardStates[guardCoord], guardState) && guardState != "" {
			return true
		}
		previousGuardStates[guardCoord] = append(previousGuardStates[guardCoord], guardState)
		tempMap, isCompleted = moveGuard(tempMap)
	}
	return false
}

func countPossibleLoops(labMap map[Coord]string, initialGuardCoord Coord, initialGuardState string) (int, int) {
	loopCount := 0
	totalCount := 0
	labMap[initialGuardCoord] = initialGuardState

	for coord, value := range labMap {
		if value == "X" {
			totalCount++
			tempMap := make(map[Coord]string)
			for k, v := range labMap {
				tempMap[k] = v
			}
			tempMap[initialGuardCoord] = initialGuardState
			if hasLoop(tempMap, coord) {
				loopCount++
			}
		}
	}
	return loopCount, totalCount
}

func main() {
	lines := readInput("input.txt")
	labMap := getLabMap(lines)
	initialGuardCoord, initialGuardState := getGuardPosition(labMap)

	tempMap := make(map[Coord]string)
	for k, v := range labMap {
		tempMap[k] = v
	}

	isCompleted := false
	for !isCompleted {
		tempMap, isCompleted = moveGuard(tempMap)
	}

	visitedCount := countVisitedLocations(tempMap)
	fmt.Println("n_locations_visited:", visitedCount)

	loopCount, totalCount := countPossibleLoops(labMap, initialGuardCoord, initialGuardState)
	fmt.Println("n_locations_with_possible_loops:", loopCount)
	fmt.Println(totalCount)
}
