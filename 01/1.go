package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func read_input() ([]int, []int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var arr1 []int
	var arr2 []int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) != 2 {
			log.Printf("wrong input for line: %s", line)
			continue
		}

		num1, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("could not convert %s to %v", parts[0], err)
			continue
		}

		num2, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("could not convert %s to %v", parts[1], err)
			continue
		}

		arr1 = append(arr1, num1)
		arr2 = append(arr2, num2)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return arr1, arr2
}

func absolute_difference(a int, b int) int {
	// if a > b then a - b will be positive
	if a > b {
		return a - b
	}

	// if b => a then b - a will be posititve
	return b - a
}

func array_total_distance(arr1 []int, arr2 []int) int {
	sort.Ints(arr1)
	sort.Ints(arr2)

	if len(arr1) != len(arr2) {
		log.Panic("arrays not of the same length")
	}

	total_distance := 0

	for i := 0; i < len(arr1); i++ {
		distance := absolute_difference(arr1[i], arr2[i])

		total_distance += distance
	}

	return total_distance
}

func array_total_similarity(arr1 []int, arr2 []int) int {
	if len(arr1) != len(arr2) {
		log.Panic("arrays not of the same length")
	}

	counter := make(map[int]int)

	for i := 0; i < len(arr1); i++ {
		counter[arr2[i]] += 1
	}

	total_similarity := 0

	for i := range arr1 {
		key := arr1[i]

		r_value := counter[key]

		similarity := key * r_value
		total_similarity += similarity
	}

	return total_similarity
}

func main() {
	arr1, arr2 := read_input()

	total_distance := array_total_distance(arr1, arr2)
	fmt.Println("distance: ", total_distance)

	total_similarity := array_total_similarity(arr1, arr2)
	fmt.Println("similarity: ", total_similarity)
}
