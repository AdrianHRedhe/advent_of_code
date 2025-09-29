package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func read_input() ([]string, []string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var towels []string
	var patterns []string

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if strings.Contains(line, ",") {
			towels = strings.Split(line, ", ")
		} else {
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return towels, patterns
}

func isValidCombination(towels []string, pattern string) (isValid bool) {
	queue := []string{}
	seen := make(map[string]bool)
	queue = append(queue, pattern)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == "" {
			return true
		}

		for _, towel := range towels {
			if len(current) < len(towel) {
				continue
			}

			if current[:len(towel)] == towel {
				subpattern := current[len(towel):]
				if !seen[subpattern] {
					seen[subpattern] = true
					queue = append(queue, subpattern)
				}
			}
		}

	}
	return false
}

func countValidCombos(towels []string, patterns []string) int {
	count := 0
	for _, pattern := range patterns {
		if isValidCombination(towels, pattern) {
			count++
		}
	}
	return count
}

func allValidCombinations(towels []string, pattern string, memory map[string]int) (nValid int) {
	// lets do recursion E.g. base case empty string
	if pattern == "" {
		return 1
	}

	if count, exists := memory[pattern]; exists {
		return count
	}

	count := 0

	for _, towel := range towels {
		if len(pattern) < len(towel) {
			continue
		}

		if pattern[:len(towel)] == towel {
			subpattern := pattern[len(towel):]
			count += allValidCombinations(towels, subpattern, memory)
		}
	}

	memory[pattern] = count
	return count
}

func countAllValidCombos(towels []string, patterns []string) int {
	count := 0
	memory := make(map[string]int)
	for _, pattern := range patterns {
		count += allValidCombinations(towels, pattern, memory)
	}
	return count
}

func main() {
	towels, patterns := read_input()
	nValidCombos := countValidCombos(towels, patterns)

	fmt.Println("part 1: ", nValidCombos)

	allValidCombos := countAllValidCombos(towels, patterns)
	fmt.Println("part 2: ", allValidCombos)
}
