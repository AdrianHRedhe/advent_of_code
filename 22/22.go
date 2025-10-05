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

func getPricesAndDeltas(secretNumber int) (prices []int, deltas []int) {
	depth := 2000
	// we need to go one level deeper than depth since the first number needs
	// to be accounted for as well
	prices = make([]int, depth+1)
	deltas = make([]int, depth+1)
	current := secretNumber

	for i := range depth + 1 {
		currentPrice := current % 10
		prices[i] = currentPrice
		if i == 0 {
			deltas[i] = 0
		} else {
			deltas[i] = prices[i] - prices[i-1]
		}
		current = nextSecretNumber(current)
	}
	return prices, deltas
}

type sequence [4]int

func calculateBananaMap(deltas []int, prices []int) (localBananasPerSeq map[sequence]int) {
	localBananasPerSeq = map[sequence]int{}
	for i := range deltas {
		// cannot start this loop before j > 3 as we would look at negative positions in list
		if i < 4 {
			continue
		}
		seq := sequence{
			deltas[i-3],
			deltas[i-2],
			deltas[i-1],
			deltas[i],
		}
		// record first position of each seq and the price at that point
		if _, exists := localBananasPerSeq[seq]; !exists {
			localBananasPerSeq[seq] = prices[i]
		}
	}
	return localBananasPerSeq
}

func calculateMaximumBananas(secretNumbers []int) (maximumBananas int) {
	globalBananasPerSeq := map[sequence]int{}

	for _, secretNumber := range secretNumbers {
		prices, deltas := getPricesAndDeltas(secretNumber)
		localBananasPerSeq := calculateBananaMap(deltas, prices)
		for seq, bananas := range localBananasPerSeq {
			globalBananasPerSeq[seq] += bananas
		}
	}

	// try out all sequences
	for _, bananas := range globalBananasPerSeq {
		if maximumBananas < bananas {
			maximumBananas = bananas
		}
	}
	return maximumBananas
}

func main() {
	secretNumbers := read_input()
	secretNumberSum := 0

	for _, secretNumber := range secretNumbers {
		secretNumberSum += nextSecretNumbers(secretNumber, 2000)
	}

	fmt.Println("part 1: ", secretNumberSum)

	maximumBananas := calculateMaximumBananas(secretNumbers)
	fmt.Println("part 2: ", maximumBananas)
}
