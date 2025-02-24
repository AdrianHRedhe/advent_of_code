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

	matrix := make([][]string, m)
	for i, row := range rows {
		matrix[i] = make([]string, n)

		for j, char := range row {
			matrix[j][i] = string(char)
		}
	}

	for _, column := range matrix {
		column := strings.Join(column, "")
		columns = append(columns, column)
	}
	return columns
}

func count_vertical(rows []string, substring string) int {
	counter := 0

	columns := make([][]string, len(rows))
	for i := range columns {
		columns[i] = make([]string, len(rows[0]))
	}

	for i, row := range rows {
		for j, char := range row {
			columns[j][i] = string(char)
		}
	}

	for _, column := range columns {
		column_as_string := strings.Join(column, "")
		counter += strings.Count(column_as_string, substring)
	}
	return counter
}

func count_left_diagonal(rows []string, substring string) int {
	counter := 0

	columns := make([][]string, len(rows))
	for i := range columns {
		columns[i] = make([]string, len(rows[0]))
	}
	return counter
}

func count_right_diagonal(rows []string, substring string) int {
	return 0
}

func count_diagonal(rows []string, substring string) int {
	counter := 0
	counter += count_left_diagonal(rows, "XMAS")
	counter += count_right_diagonal(rows, "XMAS")
	return 0
}

func count_backwards(rows []string, substring string) int {
	counter := 0
	counter += count_horizontal(rows, "XMAS")
	counter += count_vertical(rows, "XMAS")
	counter += count_diagonal(rows, "XMAS")
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
