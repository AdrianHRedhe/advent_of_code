package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func apply_rules(stone int) (stones []int) {
	// Rule 1: 0 => 1
	if stone == 0 {
		stones = append(stones, 1)
		return stones
	}

	// Rule 2: even number of digits => split into two
	stone_as_string := strconv.Itoa(stone)
	if len(stone_as_string)%2 == 0 {
		split_index := len(stone_as_string) / 2
		stone_string_1 := stone_as_string[:split_index]
		stone_string_2 := stone_as_string[split_index:]
		stone_1, err := strconv.Atoi(stone_string_1)
		if err != nil {
			log.Fatal(err)
		}
		stone_2, err := strconv.Atoi(stone_string_2)
		if err != nil {
			log.Fatal(err)
		}
		stones = append(stones, stone_1)
		stones = append(stones, stone_2)
		return stones
	}

	// Rule 3: Stone multiplied by 2024
	new_stone := stone * 2024
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
