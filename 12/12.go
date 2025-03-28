package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

func (this coord) move(other coord) coord {
	return coord{
		x: this.x + other.x,
		y: this.y + other.y,
	}
}

type plot struct {
	xy      coord
	variety string
}

type region struct {
	coords    []coord
	variety   string
	perimiter int
}

const (
	up    = 0
	down  = 1
	left  = 2
	right = 3
)

func getMatrix(lines []string) map[coord]string {
	matrix := make(map[coord]string)

	for x, line := range lines {
		chars := strings.Split(line, "")
		for y, char := range chars {
			matrix[coord{x, y}] = char
		}
	}
	return matrix
}

func findNeighbours(current plot, matrix map[coord]string) (next []coord) {
	directions := map[int]coord{
		up:    {0, 1},
		down:  {0, -1},
		left:  {-1, 0},
		right: {1, 0},
	}

	xy := current.xy
	variety := current.variety
	for _, direction := range directions {
		next_xy := xy.move(direction)
		if matrix[next_xy] == variety {
			next = append(next, next_xy)
		}
	}
	return next
}

func findRegion(current plot, matrix map[coord]string, assignedPlots []coord) (region, []coord) {
	var coordsOfRegion []coord
	var toCheck []coord
	totalPerimiter := 0

	coordsOfRegion = append(coordsOfRegion, current.xy)
	assignedPlots = append(assignedPlots, current.xy)
	toCheck = append(toCheck, current.xy)

	for len(toCheck) > 0 {
		current_xy := toCheck[0]
		next := findNeighbours(plot{current_xy, current.variety}, matrix)

		totalPerimiter += (4 - len(next))
		for idx := range next {
			switch {
			case slices.Contains(assignedPlots, next[idx]):
				continue
			case slices.Contains(toCheck, next[idx]):
				continue
			default:
				coordsOfRegion = append(coordsOfRegion, next[idx])
				assignedPlots = append(assignedPlots, next[idx])
				toCheck = append(toCheck, next[idx])
			}
		}
		toCheck = toCheck[1:]
	}
	region := region{coordsOfRegion, current.variety, totalPerimiter}
	return region, assignedPlots
}

func findRegions(matrix map[coord]string) (regions []region) {
	assignedPlots := make([]coord, 0)

	for coord, variety := range matrix {
		if slices.Contains(assignedPlots, coord) {
			continue
		}
		var region region
		region, assignedPlots = findRegion(plot{coord, variety}, matrix, assignedPlots)
		regions = append(regions, region)
	}
	return regions
}

func calculateCosts(regions []region) (totalCost int) {
	for _, region := range regions {
		area := len(region.coords)
		perimiter := region.perimiter
		totalCost += area * perimiter
	}
	return totalCost
}

func main() {
	lines := readInput()
	matrix := getMatrix(lines)
	regions := findRegions(matrix)
	totalCost := calculateCosts(regions)
	fmt.Println("total cost of fence part 1: ", totalCost)

}
