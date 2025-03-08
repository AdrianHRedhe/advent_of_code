package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() ([]int, [][]int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var test_values []int
	var variables [][]int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")

		test_value, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("could not convert %s to %v", parts[0], err)
			continue
		}

		var line_variables []int
		for _, variable := range strings.Fields(parts[1]) {
			variable_int, err := strconv.Atoi(variable)
			if err != nil {
				log.Printf("could not convert %s to %v", variable, err)
				continue
			}
			line_variables = append(line_variables, variable_int)
		}

		test_values = append(test_values, test_value)
		variables = append(variables, line_variables)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return test_values, variables
}

func operation_times(line_variables []int) []int {
	if len(line_variables) == 1 {
		return line_variables
	}
	output := make([]int, len(line_variables))
	copy(output, line_variables)
	output = output[1:]

	x, y := line_variables[0], line_variables[1]
	z := x * y
	output[0] = z
	return output
}

func operation_add(line_variables []int) []int {
	if len(line_variables) == 1 {
		return line_variables
	}
	output := make([]int, len(line_variables))
	copy(output, line_variables)
	output = output[1:]

	x, y := line_variables[0], line_variables[1]
	z := x + y
	output[0] = z
	return output
}

func action(inputs [][]int) (outputs [][]int) {
	for _, input := range inputs {
		outputs = append(outputs, operation_times(input))
		outputs = append(outputs, operation_add(input))
	}
	return outputs
}

func check_line(test_value int, line_variables []int) int {
	results := make([][]int, 1)
	results[0] = line_variables
	for range line_variables {
		results = action(results)
	}

	for _, result := range results {
		if result[0] == test_value {
			return test_value
		}
	}
	return 0
}

func main() {
	test_values, variables := read_input()

	sum_of_valid_test_values := 0
	for i := range test_values {
		test_value := test_values[i]
		line_variables := variables[i]
		sum_of_valid_test_values += check_line(test_value, line_variables)
	}
	fmt.Println("sum_of_valid_test_values:", sum_of_valid_test_values)
}
