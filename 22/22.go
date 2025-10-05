package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func read_input() []int {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var secret_numbers []int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		secret_number, err := strconv.Atoi(line)
		if err != nil {
			log.Println("Could not convert secret number", line)
			continue
		}

		secret_numbers = append(secret_numbers, secret_number)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return secret_numbers
}

var memo = map[int]int{}

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

func mix(current int, mixer int) (next int) {
	next = current ^ mixer
	return next
}

func prune(current int) (next int) {
	modNumber := 16777216
	next = modulo(current, modNumber)
	return next
}

func nextSecretNumber(current int) (next int) {
	if next, exists := memo[current]; exists {
		return next
	}
	// part 1
	mixer := current * 64
	next = mix(current, mixer)
	next = prune(next)
	// part 2
	mixer = next / 32
	next = mix(next, mixer)
	next = prune(next)
	// part 3
	mixer = next * 2048
	next = mix(next, mixer)
	next = prune(next)

	memo[current] = next
	return next
}

func nextSecretNumbers(current int, depth int) (next int) {
	for range depth {
		next = nextSecretNumber(current)
		current = next
	}
	return next
}

func main() {
	secretNumbers := read_input()
	secretNumberSum := 0

	for _, secretNumber := range secretNumbers {
		secretNumberSum += nextSecretNumbers(secretNumber, 2000)
	}

	fmt.Println("part 1: ", secretNumberSum)

	// fmt.Println("part 2: ", )
}
