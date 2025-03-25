package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func read_input() []int {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var text string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		text += line
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var stones []int
	stones_string := strings.Fields(text)
	for _, stone_string := range stones_string {
		stone, err := strconv.Atoi(stone_string)
		if err != nil {
			log.Fatal(err)
		}
		stones = append(stones, stone)
	}

	return stones
}

func split_stone(stone int) (int, int, bool) {
	// splitting stone without relying on string conversion
	// (slightly quicker but less readable.)
	// original solution converted to string and then split
	// string if the string was divisible by 2.
	var was_split bool
	n_digits := 0
	residual := stone

	for residual > 0 {
		n_digits++
		residual /= 10
	}
	// If odd return nothing
	if n_digits%2 == 1 {
		was_split = false
		return 0, 0, was_split
	}
	// If even split_stone
	divisor := 1
	for range n_digits / 2 {
		divisor *= 10
	}
	stone_left := stone / divisor
	stone_right := stone % divisor

	was_split = true
	return stone_left, stone_right, was_split
}

func apply_rules(stone int) (stones []int) {
	const multiplicator = 2024
	// Rule 1: 0 => 1
	if stone == 0 {
		stones = append(stones, 1)
		return stones
	}

	// // Rule 2: even number of digits => split into two
	stone_left, stone_right, was_split := split_stone(stone)
	if was_split {
		stones = append(stones, stone_left)
		stones = append(stones, stone_right)
		return stones
	}

	// Rule 3: Stone multiplied by 2024
	new_stone := stone * multiplicator
	stones = append(stones, new_stone)

	return stones
}

func blink(stones []int) []int {
	// Original naive implementation of blink worked ok for
	// 25 iterations
	var next_stones []int
	for _, stone := range stones {
		next_stones = append(next_stones, apply_rules(stone)...)
	}
	return next_stones
}

func blink_and_conquer(stones []int) []int {
	// The idea was to create a version which splits the
	// blink up. However, I then realized that simply using
	// memorization was a better idea.
	n_workers := 16
	next_stones := make([]int, 0, len(stones)*2)
	var wg sync.WaitGroup
	output_portal := make(chan []int, n_workers)

	input_portal := make(chan int, len(stones))
	for _, stone := range stones {
		input_portal <- stone
	}
	close(input_portal)

	for range n_workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for stone := range input_portal {
				output_portal <- apply_rules(stone)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output_portal)
	}()

	for processed_stones := range output_portal {
		next_stones = append(next_stones, processed_stones...)
	}
	return next_stones
}

func blink_and_remember(stones []int, memory map[int][]int) ([]int, map[int][]int) {
	// This was the next version of blink which remembered
	// what calculations had already been done.
	var next_stones []int
	for _, stone := range stones {
		if precalculated_stones, in_memory := memory[stone]; in_memory {
			next_stones = append(next_stones, precalculated_stones...)
		} else {
			calculated_stones := apply_rules(stone)
			next_stones = append(next_stones, calculated_stones...)
			memory[stone] = calculated_stones
		}
	}
	return next_stones, memory
}

func blink_and_remember_count_once(stone_map map[int]int, memory map[int][]int) (map[int]int, map[int][]int) {
	// Final version of blink, using memory to store
	// precalculated stones and keeping track of how many of
	// each stone exists to only count each group of stones
	// once.
	next_stone_map := make(map[int]int)
	for stone, stone_count := range stone_map {
		if precalculated_stones, in_memory := memory[stone]; in_memory {
			for _, stone := range precalculated_stones {
				next_stone_map[stone] += stone_count
			}
		} else {
			calculated_stones := apply_rules(stone)
			for _, stone := range calculated_stones {
				next_stone_map[stone] += stone_count
			}
			memory[stone] = calculated_stones
		}
	}
	return next_stone_map, memory
}

func main() {
	stones := read_input()

	// This is used to keep track of how many of each stone
	// we have rather than keeping each stone as an element
	// in a slice.
	stone_map := make(map[int]int)
	for _, stone := range stones {
		stone_map[stone] += 1
	}
	// This keeps track of all previous stone calculations.
	memory := make(map[int][]int)
	for range 25 {
		stone_map, memory = blink_and_remember_count_once(stone_map, memory)
	}
	count := 0
	for _, stone_count := range stone_map {
		count += stone_count
	}
	fmt.Println("After blinking 25 times part 1:", count)
	for range 50 {
		stone_map, memory = blink_and_remember_count_once(stone_map, memory)
	}
	count = 0
	for _, stone_count := range stone_map {
		count += stone_count
	}
	fmt.Println("After blinking 75 times part 2:", count)
}
