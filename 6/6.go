package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

func get_lab_map(lines []string) map[coord]string {
	lab_map := make(map[coord]string)

	for x, line := range lines {
		chars := strings.Split(line, "")
		for y, char := range chars {
			lab_map[coord{x, y}] = char
		}
	}

	return lab_map
}

func action(lab_map map[coord]string) (updated_lab_map map[coord]string, completed bool) {
	guard_states := []string{"^", ">", "v", "<"}
	movement_states := []coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	// a) find the tracker and the direction
	for xy, value := range lab_map {
		if slices.Contains(guard_states, value) {
			guard_state_index := slices.Index(guard_states, value)
			movement_state := movement_states[guard_state_index]

			x, y := xy.x, xy.y
			dx, dy := movement_state.x, movement_state.y

			// b) if next index is outside the board then replace
			// with X and return has has_changed
			next_value, is_inside_map := lab_map[coord{x + dx, y + dy}]

			if !is_inside_map {
				lab_map[coord{x, y}] = "X"
				return lab_map, true
			}

			// c) if next index does have a # then only turn it and
			// call function again
			if next_value == "#" {
				next_guard_state := (guard_state_index + 1) % len(guard_states)
				lab_map[coord{x, y}] = guard_states[next_guard_state]
				return lab_map, false
			}

			// d) if next index does not have a # then move there
			// and put X on old location
			lab_map[coord{x + dx, y + dy}] = lab_map[coord{x, y}]
			lab_map[coord{x, y}] = "X"
			return lab_map, false
		}
	}
	return lab_map, false
}

func count_locations(lab_map map[coord]string) int {
	counter := 0
	for _, v := range lab_map {
		if v == "X" {
			counter++
		}
	}
	return counter
}

func main() {
	lines := read_input()
	// part 1:
	lab_map := get_lab_map(lines)
	is_completed := false
	for !is_completed {
		lab_map, is_completed = action(lab_map)
	}
	n_locations_visited := count_locations(lab_map)
	// total_similarity := array_total_similarity(arr1, arr2)
	fmt.Println("n_locations_visited: ", n_locations_visited)
}
