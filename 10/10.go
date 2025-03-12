package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() []string {
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

type coord_3d struct {
	xy coord
	z  int
}

const (
	up    = 0
	down  = 1
	left  = 2
	right = 3
)

func get_matrix(lines []string) map[coord]int {
	matrix := make(map[coord]int)

	for x, line := range lines {
		chars := strings.Split(line, "")
		for y, char := range chars {
			height, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatal("char couldnt be converted to int")
			}
			matrix[coord{x, y}] = height
		}
	}
	return matrix
}

func check_trailhead_count_based_on_score(xy coord, matrix map[coord]int) int {
	var possible_locations []coord_3d
	directions := map[int]coord{
		up:    {0, 1},
		down:  {0, -1},
		left:  {-1, 0},
		right: {1, 0},
	}
	for _, direction := range directions {
		if matrix[xy.move(direction)] == 1 {
			location := coord_3d{xy: xy.move(direction), z: 1}
			possible_locations = append(possible_locations, location)
		}
	}
	is_complete := false
	for !is_complete {
		possible_locations, is_complete = update_possible_locations(possible_locations, directions, matrix)
	}
	// Only count each finish line once
	number_of_valid_trails := 0
	var visited_high_points []coord
	for _, loc := range possible_locations {
		has_visited := false
		for _, prev_loc := range visited_high_points {
			if (loc.xy.x == prev_loc.x) && (loc.xy.y == prev_loc.y) {
				has_visited = true
			}
		}
		if !has_visited {
			number_of_valid_trails++
		}
		visited_high_points = append(visited_high_points, loc.xy)
	}
	return number_of_valid_trails
}

func update_possible_locations(possible_locations []coord_3d, directions map[int]coord, matrix map[coord]int) ([]coord_3d, bool) {
	var next_possible_locations []coord_3d
	if len(possible_locations) == 0 {
		return possible_locations, true
	}

	any_non_nine_height := false
	for _, possible_location := range possible_locations {
		xy, height := possible_location.xy, possible_location.z
		if height != 9 {
			any_non_nine_height = true
		}
		for _, direction := range directions {
			if matrix[xy.move(direction)] == height+1 {
				next_valid_location := coord_3d{xy: xy.move(direction), z: height + 1}
				next_possible_locations = append(next_possible_locations, next_valid_location)
			}
		}
	}

	if !any_non_nine_height {
		return possible_locations, true
	}

	return next_possible_locations, false
}

func check_trailhead_count_based_on_rating(xy coord, matrix map[coord]int) int {
	var possible_locations []coord_3d
	directions := map[int]coord{
		up:    {0, 1},
		down:  {0, -1},
		left:  {-1, 0},
		right: {1, 0},
	}
	for _, direction := range directions {
		if matrix[xy.move(direction)] == 1 {
			location := coord_3d{xy: xy.move(direction), z: 1}
			possible_locations = append(possible_locations, location)
		}
	}
	is_complete := false
	for !is_complete {
		possible_locations, is_complete = update_possible_locations(possible_locations, directions, matrix)
	}
	// Only count each finish line once
	return len(possible_locations)
}

func main() {
	lines := read_input()
	matrix := get_matrix(lines)

	trail_counter_part_1 := 0
	for xy := range matrix {
		if matrix[xy] == 0 {
			trail_counter_part_1 += check_trailhead_count_based_on_score(xy, matrix)
		}
	}
	println("trail_counter_part_1:", trail_counter_part_1)

	trail_counter_part_2 := 0
	for xy := range matrix {
		if matrix[xy] == 0 {
			trail_counter_part_2 += check_trailhead_count_based_on_rating(xy, matrix)
		}
	}
	println("trail_counter_part_2:", trail_counter_part_2)
}
