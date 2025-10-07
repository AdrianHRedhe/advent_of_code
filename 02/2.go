package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() [][]int {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var reports [][]int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		var report []int

		for i := range parts {
			layer, err := strconv.Atoi(parts[i])

			if err != nil {
				log.Printf("could not convert %s to %v", parts[i], err)
				continue
			}

			report = append(report, layer)
		}

		reports = append(reports, report)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return reports
}

func absolute_difference(a int, b int) int {
	// if a > b then a - b will be positive
	if a > b {
		return a - b
	}

	// if b => a then b - a will be posititve
	return b - a
}

func check_report(report []int) bool {
	var differences []int
	for i := 0; i < len(report)-1; i++ {
		current := report[i]
		next := report[i+1]
		diff := next - current

		differences = append(differences, diff)
	}

	// Magnitude of difference is between 1-3
	for _, diff := range differences {
		if diff < 0 {
			diff = -diff
		}
		if (diff < 1) || (3 < diff) {
			return false
		}

	}

	// All differences are either pos / neg
	for i := 0; i < len(differences)-1; i++ {
		current := differences[i]
		next := differences[i+1]

		if (current < 0) && (next > 0) {
			return false
		}

		if (current > 0) && (next < 0) {
			return false
		}
	}

	return true

}

func check_report_with_dampener(report []int) bool {
	var differences []int
	for i := 0; i < len(report)-1; i++ {
		current := report[i]
		next := report[i+1]
		diff := next - current

		differences = append(differences, diff)
	}

	var bad_layer int
	first_strike := false
	// Magnitude of difference is between 1-3
	for i, diff := range differences {
		if diff < 0 {
			diff = -diff
		}

		if (diff < 1) || (3 < diff) {
			if !first_strike {
				first_strike = true
				bad_layer = i
			} else {
				return false
			}
		}
	}

	first_strike = false
	// All differences are either pos / neg
	for i := 0; i < len(differences)-1; i++ {
		current := differences[i]
		next := differences[i+1]

		if (current < 0) && (next > 0) {
			if !first_strike {
				first_strike = true

				if bad_layer != i {
					return false
				}
			} else {
				return false
			}
		}

		if (current > 0) && (next < 0) {
			if !first_strike {
				first_strike = true

				if bad_layer != i {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func check_reports(reports [][]int) (int, int) {

	safe_report_counter := 0
	safe_report_with_dampener_counter := 0

	for i := range reports {
		report := reports[i]
		is_safe := check_report(report)
		is_safe_with_dampener := check_report_with_dampener(report)

		if is_safe {
			safe_report_counter += 1
		}

		if is_safe_with_dampener {
			safe_report_with_dampener_counter += 1
		}
	}

	return safe_report_counter, safe_report_with_dampener_counter
}

func main() {
	reports := read_input()

	n_safe_reports, n_safe_reports_with_dampener := check_reports(reports)
	fmt.Println("safe report counter: ", n_safe_reports)
	fmt.Println("safe report with dampener counter: ", n_safe_reports_with_dampener)

}
