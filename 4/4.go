package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//"regexp"
	"strings"
)

func read_input() []string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string

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

//r := regexp.MustCompile("XMAS")
//n := r.FindAllString(lines, -1)

// return len(n) //strings.Count(lines, "xmas")

func count_horizontal(rows []string, substring string) int {
	counter := 0
	// check if possible with continuing outside of lines
	for _, row := range rows {
		counter += strings.Count(row, substring)
	}
	return counter
}

func transpose(rows []string) []string {
	n, m := len(rows), len(rows[0])

	columns := make([]string, m)

	matrix := make([][]rune, m)
	for j := range m {
		matrix[j] = make([]rune, n)
	}

	for i, row := range rows {
		for j, char := range row {
			matrix[j][i] = char
		}
	}

	for i := range matrix {
		columns[i] = string(matrix[i])
	}
	return columns
}

func count_vertical(rows []string, substring string) int {
	columns := transpose(rows)
	return count_horizontal(columns, substring)
}

func count_left_diagonal(rows []string, substring string) int {
	counter := 0
	n, m := len(rows), len(rows[0])
	matrix := make([][]string, n)
	for i := range matrix {
		matrix[i] = strings.Split(rows[i], "")
	}

	for i := range matrix {
		if i >= (n - len(substring)) {
			continue
		}
		for j := range matrix {
			if j >= (m - len(substring)) {
				continue
			}
			possible_substring := ""
			for idx := range len(substring) {
				possible_substring += matrix[i+idx][j+idx]
			}
			if possible_substring == substring {
				counter += 1
			}
		}
	}

	return counter
}

func count_right_diagonal(rows []string, substring string) int {
	// columns := transpose(rows)
	// fmt.Println(len(columns), len(columns[0]))
	return 0 //count_left_diagonal(columns, substring)
}

func count_diagonal(rows []string, substring string) int {
	counter := 0
	counter += count_left_diagonal(rows, substring)
	counter += count_right_diagonal(rows, substring)
	return counter
}

func count_backwards(rows []string, substring string) int {
	string_as_runes := []rune(substring)
	reversed_string_as_runes := make([]rune, len(substring))

	for i, char := range string_as_runes {
		reversed_string_as_runes[len(substring)-i-1] = char
	}

	reversed_substring := string(reversed_string_as_runes)
	fmt.Println(reversed_substring)
	counter := 0
	counter += count_horizontal(rows, reversed_substring)
	counter += count_vertical(rows, reversed_substring)
	counter += count_diagonal(rows, reversed_substring)
	return 0
}

func count_xmas(rows []string) int {
	n_xmas := 0
	n_xmas += count_horizontal(rows, "XMAS")
	n_xmas += count_vertical(rows, "XMAS")
	n_xmas += count_diagonal(rows, "XMAS")
	n_xmas += count_backwards(rows, "XMAS")

	return n_xmas
}

func main() {
	rows := read_input()

	n_xmas := count_xmas(rows)
	fmt.Println("XMAS counter 1: ", n_xmas)

	// fmt.Println("XMAS counter 1: ", total_distance)
}
