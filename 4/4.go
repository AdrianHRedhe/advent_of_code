package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func read_input() []string {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lines []string

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

type coord struct {
	x, y int
}

func get_matrix(lines []string) map[coord]string {
	matrix := make(map[coord]string)

	for x, line := range lines {
		chars := strings.Split(line, "")
		for y, char := range chars {
			matrix[coord{x, y}] = char
		}
	}
	return matrix
}

func check_all_directions(matrix map[coord]string, start coord, substring string) int {
	counter := 0
	x, y := start.x, start.y

	for _, dx := range []int{1, 0, -1} {
		for _, dy := range []int{1, 0, -1} {
			if !(dx == 0 && dy == 0) {
				// Possible direction
				found_string := ""
				for i := range substring {
					found_string += matrix[coord{x + dx*i, y + dy*i}]
				}
				if found_string == substring {
					counter += 1
				}
			}
		}
	}
	return counter

}

func count_all_matches(rows []string, substring string) int {
	n_matches := 0
	matrix := get_matrix(rows)

	for xy_pair := range matrix {
		n_matches += check_all_directions(matrix, xy_pair, substring)
	}

	return n_matches
}

func x_mas(rows []string) int {
	counter := 0
	matrix := get_matrix(rows)

	for xy_pair := range matrix {
		counter += check_x_mas(matrix, xy_pair)
	}
	return counter
}

func check_x_mas(matrix map[coord]string, xy coord) int {
	nw := coord{xy.x + 1, xy.y + 1}
	ne := coord{xy.x + 1, xy.y - 1}
	se := coord{xy.x - 1, xy.y - 1}
	sw := coord{xy.x - 1, xy.y + 1}

	if matrix[xy] != "A" {
		return 0
	}

	if !((matrix[nw] == "M" && matrix[se] == "S") || (matrix[nw] == "S" && matrix[se] == "M")) {
		return 0
	}

	if !((matrix[ne] == "M" && matrix[sw] == "S") || (matrix[ne] == "S" && matrix[sw] == "M")) {
		return 0
	}
	return 1
}

func main() {
	rows := read_input()
	n_xmas := count_all_matches(rows, "XMAS")
	// n_xmas := count_xmas(rows)
	fmt.Println("XMAS counter 1: ", n_xmas)

	x_mas_counter := x_mas(rows)
	fmt.Println("X-MAS counter 2: ", x_mas_counter)
}
