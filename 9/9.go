package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func read_input() string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var text string

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

func decompress(text string) []int {
	value_denotes_free_space := false
	id := 0
	decompressed_ints := []int{}

	for _, char := range text {
		value, err := strconv.Atoi(string(char))
		if err != nil {
			log.Fatal("could not convert to int")
		}
		if value_denotes_free_space {
			for range value {
				decompressed_ints = append(decompressed_ints, -1)
			}
			value_denotes_free_space = !value_denotes_free_space
		} else {
			for range value {
				decompressed_ints = append(decompressed_ints, id)
			}
			id++
			value_denotes_free_space = !value_denotes_free_space
		}
	}
	return decompressed_ints
}

func find_index_of_last_number(ints []int) int {
	for i := range ints {
		value := ints[len(ints)-i-1]
		if value != -1 {
			return len(ints) - i - 1
		}
	}
	return 0
}

func recompress(ints []int) []int {
	for i := range ints {
		if ints[i] == -1 {
			index_of_last_number := find_index_of_last_number(ints)
			if index_of_last_number > i {
				ints[i], ints[index_of_last_number] = ints[index_of_last_number], ints[i]
			} else {
				break
			}
		}
	}
	return ints
}

func checksum(ints []int) int {
	checksum := 0
	for i, value := range ints {
		if value == -1 {
			continue
		}
		checksum += i * value
	}
	return checksum
}

func try_to_recompress_file(id int, ints []int) []int {
	file_size := 0
	var first_index_of_file int
	for i, value := range ints {
		if value == id {
			if file_size == 0 {
				first_index_of_file = i
			}
			file_size++
		}
	}

	for i := range ints {
		if ints[i] == -1 {
			contigious_space := 1

			if i >= first_index_of_file {
				// We have reached the location where the
				// file is, it cannot be more compressed.
				return ints
			}

			if i+1 >= len(ints) {
				// We are in the last spot so no changes
				// were made to the array so we can return
				return ints
			}

			for ints[i+contigious_space] == -1 {
				if contigious_space+i >= len(ints)-1 {
					break
				}
				contigious_space++
			}

			if file_size > contigious_space {
				continue
			}
			// Space is large enough to fit file
			for j := range file_size {
				ints[i+j], ints[first_index_of_file+j] = ints[first_index_of_file+j], ints[i+j]
			}
			// File has been moved so we can now go back.
			return ints
		}
	}
	// No spaces fond
	return ints
}

func recompress_whole_files(ints []int) []int {
	var highest_id int
	for i := range ints {
		if ints[len(ints)-i-1] >= -1 {
			highest_id = ints[len(ints)-i-1]
			break
		}
	}

	for i := range highest_id {
		cur_file_id := highest_id - i
		ints = try_to_recompress_file(cur_file_id, ints)
	}
	return ints
}

func main() {
	text := read_input()
	decompressed_ints := decompress(text)
	recompressed_ints := recompress(decompressed_ints)
	checksum_part1 := checksum(recompressed_ints)
	println("checksum_part1:", checksum_part1)

	decompressed_ints = decompress(text)
	recompressed_files := recompress_whole_files(decompressed_ints)
	checksum_part2 := checksum(recompressed_files)
	println("checksum_part2:", checksum_part2)
}
