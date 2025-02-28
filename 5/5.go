package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func read_input() ([]string, []string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var rules []string
	var inputs []string
	empty_line := false

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			empty_line = true
		}

		if !empty_line {
			rules = append(rules, line)
		} else {
			inputs = append(inputs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return rules, inputs
}

func get_rule_maps(rules []string) (map[string][]string, map[string][]string) {
	rule_map_prev := make(map[string][]string)
	rule_map_next := make(map[string][]string)

	for _, rule := range rules {
		split_rule := strings.Split(rule, "|")
		prev, next := split_rule[0], split_rule[1]
		rule_map_next[prev] = append(rule_map_next[prev], next)
		rule_map_prev[next] = append(rule_map_prev[next], prev)
	}

	return rule_map_next, rule_map_prev
}

func check_one_idx(idx int, instructions []string, rule_map_next map[string][]string, rule_map_prev map[string][]string) bool {
	current := instructions[idx]
	var previous []string
	var upcoming []string

	if idx == 0 {
		previous = []string{}
	} else {
		previous = instructions[0 : idx-1]
	}

	if idx == len(instructions)-1 {
		upcoming = []string{}
	} else {
		upcoming = instructions[idx+1:]
	}

	for _, prev := range previous {
		if slices.Contains(rule_map_prev[prev], current) {
			// Current should be before Prev => rule is broken
			return false
		}
		if slices.Contains(rule_map_next[current], prev) {
			// Prev should be after Current => rule is broken
			return false
		}
	}

	for _, next := range upcoming {
		if slices.Contains(rule_map_prev[current], next) {
			// Next should be before Current => rule is broken
			return false
		}
		if slices.Contains(rule_map_next[next], current) {
			// Current should be after Next => rule is broken
			return false
		}
	}

	return true
}

func check_one_input(input string, rule_map_next map[string][]string, rule_map_prev map[string][]string) int {
	instructions := strings.Split(input, ",")
	for idx := range instructions {
		if !check_one_idx(idx, instructions, rule_map_next, rule_map_prev) {
			return 0
		}
	}
	middle_index := int(len(instructions) / 2)
	return_value, err := strconv.Atoi(instructions[middle_index])
	if err != nil {
		log.Fatal(err)
	}
	return return_value
}

func part_1(rules []string, inputs []string) int {
	rule_map_next, rule_map_prev := get_rule_maps(rules)
	counter := 0
	for _, input := range inputs[1:] {
		counter += check_one_input(input, rule_map_next, rule_map_prev)
	}

	return counter
}

func main() {
	rules, inputs := read_input()

	part_1 := part_1(rules, inputs)
	fmt.Println(part_1)
}
