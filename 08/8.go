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

func find_antinodes_based_on_frequency(matrix map[coord]string, frequency string, anti_nodes_locations map[coord]string) map[coord]string {
	for xy_a, value := range matrix {
		if value != frequency {
			continue
		}

		for xy_b, value := range matrix {
			// needs to be of same frequecy and another
			// point
			if (value != frequency) || (xy_a == xy_b) {
				continue
			}
			dx := xy_a.x - xy_b.x
			dy := xy_a.y - xy_b.y

			anti_node_location := coord{xy_a.x + dx, xy_a.y + dy}
			// need to be in bounds of the map and is not if
			// this value is not recorded in original matrix
			if matrix[anti_node_location] == "" {
				continue
			}
			anti_nodes_locations[anti_node_location] = "#"
		}
	}
	return anti_nodes_locations
}

func find_antinodes_based_on_frequency_with_resonance(matrix map[coord]string, frequency string, anti_nodes_locations map[coord]string) map[coord]string {
	for xy_a, value := range matrix {
		if value != frequency {
			continue
		}

		for xy_b, value := range matrix {
			// needs to be of same frequecy
			if value != frequency {
				continue
			}
			dx := xy_a.x - xy_b.x
			dy := xy_a.y - xy_b.y

			it := 0
			anti_node_location := coord{xy_a.x + dx, xy_a.y + dy}
			// need to be in bounds of the map and is not if
			// this value is not recorded in original matrix

			for matrix[anti_node_location] != "" {
				anti_nodes_locations[anti_node_location] = "#"
				it++
				anti_node_location = coord{xy_a.x + dx*(it+1), xy_a.y + dy*(it+1)}
				if dx == 0 && dy == 0 {
					// If it is the same point it only needs
					// to be added once
					break
				}
			}
		}
	}
	return anti_nodes_locations
}

func main() {
	lines := read_input()
	matrix := get_matrix(lines)

	anti_nodes_locations := make(map[coord]string)
	for _, frequency := range matrix {
		// . is not a frequency
		if frequency == "." {
			continue
		}
		anti_nodes_locations = find_antinodes_based_on_frequency(matrix, frequency, anti_nodes_locations)
	}

	counter := 0
	for range anti_nodes_locations {
		counter++
	}
	fmt.Println("unique anti-nodes:", counter)

	anti_nodes_locations_with_resonance := make(map[coord]string)
	for _, frequency := range matrix {
		// . is not a frequency
		if frequency == "." {
			continue
		}
		anti_nodes_locations_with_resonance = find_antinodes_based_on_frequency_with_resonance(matrix, frequency, anti_nodes_locations_with_resonance)
	}
	counter = 0
	for range anti_nodes_locations_with_resonance {
		counter++
	}
	fmt.Println("unique anti-nodes with resonance:", counter)
}
