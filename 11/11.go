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

	// Alternative version of Rule 2.
	// stone_as_string := strconv.Itoa(stone)
	// if len(stone_as_string)%2 == 0 {
	// 	split_index := len(stone_as_string) / 2
	// 	stone_string_1 := stone_as_string[:split_index]
	// 	stone_string_2 := stone_as_string[split_index:]
	// 	stone_1, err := strconv.Atoi(stone_string_1)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	stone_2, err := strconv.Atoi(stone_string_2)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	stones = append(stones, stone_1)
	// 	stones = append(stones, stone_2)
	// 	return stones
	// }

	// Rule 3: Stone multiplied by 2024
	new_stone := stone * multiplicator
	stones = append(stones, new_stone)

	return stones
}

func blink(stones []int) []int {
	var next_stones []int
	for _, stone := range stones {
		next_stones = append(next_stones, apply_rules(stone)...)
	}
	return next_stones
}

func blink_and_conquer(stones []int) []int {
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

func main() {
	stones := read_input()
	for range 25 {
		stones = blink(stones)
	}
	fmt.Println("After blinking 25 times part 1:", len(stones))
	it := 25
	for range 50 {
		log.Println(it)
		stones = blink(stones)
		it++
	}
	fmt.Println("After blinking 75 times part 2:", len(stones))
}
