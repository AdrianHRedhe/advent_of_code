package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readInput() (lines []string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scan and format lines.
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

type key [5]int
type lock [5]int

func (l lock) overlaps(k key) bool {
	for x := range 5 {
		if l[x]+k[x] > 5 {
			return true
		}
	}
	return false
}

func formatLock(lines []string) lock {
	formattedLock := lock{}
	for x := range 5 {
		for revY := range 6 {
			y := 5 - revY
			if string(lines[y][x]) == "#" {
				formattedLock[x] = y
				break
			}
		}
	}
	return formattedLock
}

func formatKey(lines []string) key {
	formattedKey := key{}
	for x := range 5 {
		for y := range 6 {
			if string(lines[y][x]) == "#" {
				formattedKey[x] = 5 - y
				break
			}
		}
	}
	return formattedKey
}

func formatInput(lines []string) (locks []lock, keys []key) {
	for i := 0; i <= len(lines); i += 8 {
		if lines[i] == "#####" {
			locks = append(locks, formatLock(lines[i:i+6]))
		} else {
			keys = append(keys, formatKey(lines[i+1:i+7]))
		}
	}
	return locks, keys
}

func part1(locks []lock, keys []key) int {
	count := 0
	for _, l := range locks {
		for _, k := range keys {
			if !l.overlaps(k) {
				count++
			}
		}
	}
	return count
}

func main() {
	lines := readInput()
	locks, keys := formatInput(lines)
	fmt.Println("part 1:", part1(locks, keys))
}
