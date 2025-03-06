package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"
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

func get_guard_position(lab_map map[coord]string) (coord, string) {
	guard_states := []string{"^", ">", "v", "<"}
	for xy, value := range lab_map {
		if slices.Contains(guard_states, value) {
			return xy, value
		}
	}
	return coord{-1, -1}, ""
}

func has_loop(temp_lab_map map[coord]string, location_of_obstacle coord) (has_loop bool) {
	temp_lab_map[location_of_obstacle] = "#"
	previous_guard_locations := make(map[coord][]string)

	is_completed := false
	for !is_completed {
		temp_lab_map, is_completed = action(temp_lab_map)

		guard_position, guard_state := get_guard_position(temp_lab_map)
		if slices.Contains(previous_guard_locations[guard_position], guard_state) {
			if guard_state != "" {
				return true
			}
		}

		previous_guard_locations[guard_position] = append(previous_guard_locations[guard_position], guard_state)
	}
	return false
}

func count_all_possible_loops(lab_map map[coord]string, original_guard_position coord, original_guard_state string) int {
	// The lab map has already found the possible candidates
	// to put in an obstacle but is missing the guard.
	lab_map[original_guard_position] = original_guard_state

	// Use concurrency in GO.
	var wg sync.WaitGroup
	results := make(chan int)

	for xy, value := range lab_map {
		if value == "X" {
			wg.Add(1)

			go func(xy coord) {
				defer wg.Done()
				temp_lab_map := make(map[coord]string)
				for xy, value := range lab_map {
					temp_lab_map[xy] = value
				}

				if has_loop(temp_lab_map, xy) {
					results <- 1
				} else {
					results <- 0
				}
			}(xy)
		}
	}

	go func() {
		wg.Wait()
		close(results) // Close the channel once all goroutines are done
	}()

	counter := 0
	for result := range results {
		counter += result
	}

	return counter
}

func main() {
	lines := read_input()
	// part 1:
	lab_map := get_lab_map(lines)
	original_guard_position, original_guard_state := get_guard_position(lab_map)
	is_completed := false
	for !is_completed {
		lab_map, is_completed = action(lab_map)
	}
	n_locations_visited := count_locations(lab_map)
	fmt.Println("n_locations_visited: ", n_locations_visited)

	// part 2
	n_locations_with_possible_loops := count_all_possible_loops(lab_map, original_guard_position, original_guard_state)
	fmt.Println("n_locations_with_possible_loops:", n_locations_with_possible_loops)
}
