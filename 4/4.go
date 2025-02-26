package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//"regexp"
	// "strings"
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

func make_matrix(rows []string) [][]string {
	n, m := len(rows), len(rows[0])

	matrix := make([][]string, m)
	for j := range m {
		matrix[j] = make([]string, n)
	}

	for i, row := range rows {
		for j, char := range row {
			matrix[j][i] = string(char)
		}
	}

	return matrix
}

type coord struct {
	x, y int
}

func make_matrix_map(rows []string) map[coord]string {

	// X, Y := len(rows), len(rows[0])
	matrix := make(map[coord]string)

	for i, row := range rows {
		for j, char := range row {
			matrix[coord{i, j}] = string(char)
		}
	}

	return matrix
}

func check_all_directions(matrix map[coord]string, start coord, substring string) int {
	counter := 0
	x, y := start.x, start.y

	possible_directions := []coord{{1, -1}, {1, 0}, {0, 1}, {0, -1}, {0, 1}, {-1, -1}, {-1, 0}, {-1, 1}}
	for _, possible_direction := range possible_directions {
		dx, dy := possible_direction.x, possible_direction.y
		found_string := ""
		for i := range substring {
			found_string += matrix[coord{x + dx*i, y + dy*i}]
		}

		if found_string == substring {
			counter += 1
		}
	}

	return counter
}

func count_all_matches(rows []string, substring string) int {
	n_matches := 0
	matrix := make_matrix_map(rows)

	for xy_pair := range matrix {
		n_matches += check_all_directions(matrix, xy_pair, substring)
	}

	return n_matches
}

func main() {
	rows := read_input()
	n_xmas := count_all_matches(rows, "XMAS")
	// n_xmas := count_xmas(rows)
	fmt.Println("XMAS counter 1: ", n_xmas)
}
