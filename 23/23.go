package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func readInput() []string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var connections []string

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		connection := scanner.Text()
		connections = append(connections, connection)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return connections
}

var connect2 = map[string][]string{}
var connect3 = map[string]bool{}

func updateConnect3(a, b string) {
	for _, c1 := range connect2[a] {
		for _, c2 := range connect2[b] {
			if c1 != c2 {
				continue
			}
			c := string(c1)
			if c != a && c != b {
				computers := []string{a, b, c}
				sort.Strings(computers)
				key := computers[0] + "-" + computers[1] + "-" + computers[2]
				connect3[key] = true
			}
		}

	}
	return
}

func part1(connections []string) int {

	for _, connection := range connections {
		ab := strings.Split(connection, "-")
		a, b := ab[0], ab[1]
		connect2[a] = append(connect2[a], b)
		connect2[b] = append(connect2[b], a)
		updateConnect3(a, b)
	}
	count := 0
	for key := range connect3 {
		for _, computer := range strings.Split(key, "-") {
			if string(computer[0]) == "t" {
				count++
				break
			}
		}
	}
	return count
}

func main() {
	connections := readInput()
	count := part1(connections)

	fmt.Println("part 1: ", count)
}
