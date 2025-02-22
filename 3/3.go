package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var text string

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		text += line
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return text
}

func split_parse_multiplicate(text string) int {
	parts := strings.Split(text, ")")
	total_sum := 0

	// split the instructions
	for _, part := range parts {
		total_sum += parse_and_multiplicate(part)
	}

	return total_sum
}

func split_parse_multiplicate_with_do_instuction(text string) int {
	parts := strings.Split(text, ")")
	total_sum := 0
	do := true
	var result int

	// split the instructions
	for _, part := range parts {
		result, do = parse_and_multiplicate_with_do_instructions(part, do)
		if do {
			total_sum += result
		}
	}

	return total_sum
}

func parse_and_multiplicate(part string) int {
	// if mul( not in string then return 0
	if !strings.Contains(part, "mul(") {
		return 0
	}

	contenders := strings.SplitAfter(part, "mul(")
	contender := contenders[len(contenders)-1]

	ab := strings.Split(contender, ",")
	if len(ab) != 2 {
		return 0
	}

	a, b := ab[0], ab[1]

	if (len(a) < 1) || (3 < len(a)) {
		return 0
	}
	if (len(b) < 1) || (3 < len(b)) {
		return 0
	}

	int_a, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	int_b, err := strconv.Atoi(b)
	if err != nil {
		return 0
	}

	return int_a * int_b
}

func parse_and_multiplicate_with_do_instructions(part string, do bool) (int, bool) {
	// Check if there is an update to the do instruction
	if strings.Contains(part, "do") {
		do_index := strings.Index(part, "do(")
		dont_index := strings.Index(part, "don't(")

		if do_index != dont_index {
			if do_index > dont_index {
				do = true
			} else {
				do = false
			}
		}
	}

	// if mul( not in string then return 0
	if !strings.Contains(part, "mul(") {
		return 0, do
	}

	contenders := strings.SplitAfter(part, "mul(")
	contender := contenders[len(contenders)-1]

	ab := strings.Split(contender, ",")
	if len(ab) != 2 {
		return 0, do
	}

	a, b := ab[0], ab[1]

	if (len(a) < 1) || (3 < len(a)) {
		return 0, do
	}
	if (len(b) < 1) || (3 < len(b)) {
		return 0, do
	}

	int_a, err := strconv.Atoi(a)
	if err != nil {
		return 0, do
	}
	int_b, err := strconv.Atoi(b)
	if err != nil {
		return 0, do
	}

	return int_a * int_b, do
}

func main() {
	text := read_input()

	multi := split_parse_multiplicate(text)
	fmt.Println("mutli: ", multi)

	multi = split_parse_multiplicate_with_do_instuction(text)
	fmt.Println("mutli with do instructions: ", multi)

}
