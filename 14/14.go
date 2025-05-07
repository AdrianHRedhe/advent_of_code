package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput() []string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lines []string

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

func modulo(x int, d int) int {
	// Unfortunately Golang's % operator is not modulo
	// at least not in the common sense but instead
	// remainder.
	m := x % d
	if x < 0 && d < 0 {
		m -= d
	}
	if x < 0 && d > 0 {
		m += d
	}
	return m
}

func (this coord) moveAndTeleport(other coord) coord {
	width := 101
	height := 103
	return coord{
		x: modulo((this.x + other.x), width),
		y: modulo((this.y + other.y), height),
	}
}

type guard struct {
	position coord
	velocity coord
}

type quadrant struct {
	NW coord
	SE coord
}

func stringToCoord(string string) coord {
	coordsSeparatedByComma := string[2:]
	coordStrings := strings.Split(coordsSeparatedByComma, ",")
	x_str, y_str := coordStrings[0], coordStrings[1]

	x, err := strconv.Atoi(x_str)
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.Atoi(y_str)
	if err != nil {
		log.Fatal(err)
	}

	return coord{
		x: x,
		y: y,
	}
}

func getGuards(lines []string) (guards []guard) {
	for _, line := range lines {
		inputs := strings.Fields(line)
		position := stringToCoord(inputs[0])
		velocity := stringToCoord(inputs[1])
		guard := guard{position, velocity}
		guards = append(guards, guard)
	}
	return guards
}

func moveOneRound(guards []guard) (movedGuards []guard) {
	for _, guard := range guards {
		nextPosition := guard.position.moveAndTeleport(guard.velocity)
		guard.position = nextPosition
		movedGuards = append(movedGuards, guard)
	}
	return movedGuards
}

func getQuadrants() (quadrants []quadrant) {
	// the parts of the middle dont count
	width := 101
	height := 103
	Q1 := quadrant{
		NW: coord{x: 0, y: 0},
		SE: coord{x: (width / 2) - 1, y: (height / 2) - 1},
	}
	Q2 := quadrant{
		NW: coord{x: (width / 2) + 1, y: 0},
		SE: coord{x: width - 1, y: (height / 2) - 1},
	}
	Q3 := quadrant{
		NW: coord{x: 0, y: (height / 2) + 1},
		SE: coord{x: (width / 2) - 1, y: height - 1},
	}
	Q4 := quadrant{
		NW: coord{x: (width / 2) + 1, y: (height / 2) + 1},
		SE: coord{x: width - 1, y: height - 1},
	}

	quadrants = []quadrant{Q1, Q2, Q3, Q4}
	return quadrants
}

func guardInQuadrant(guard guard, quadrant quadrant) bool {
	xMin, xMax := quadrant.NW.x, quadrant.SE.x
	if !((xMin <= guard.position.x) && (guard.position.x <= xMax)) {
		return false
	}

	yMin, yMax := quadrant.NW.y, quadrant.SE.y
	if !((yMin <= guard.position.y) && (guard.position.y <= yMax)) {
		return false
	}
	return true
}

func calculateSecurityScore(guards []guard, quadrants []quadrant) (securityScore int) {
	nGuardsPerQuadrant := []int{0, 0, 0, 0}
	for _, guard := range guards {
		for idx, quadrant := range quadrants {
			if guardInQuadrant(guard, quadrant) {
				nGuardsPerQuadrant[idx]++
			}
		}
	}
	securityScore = 1
	for _, nGuards := range nGuardsPerQuadrant {
		securityScore *= nGuards
		println(nGuards)
	}

	return securityScore
}

func main() {
	lines := readInput()
	guards := getGuards(lines)
	quadrants := getQuadrants()

	for range 100 {
		guards = moveOneRound(guards)
	}
	securityScore := calculateSecurityScore(guards, quadrants)

	fmt.Println("security score part 1: ", securityScore)

	// for part 2 it requires visual confirmation so I will
	// simply let it run a few times and look for the signs
	guards = getGuards(lines)

}
