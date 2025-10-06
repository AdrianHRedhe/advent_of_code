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

func read_input() (wires map[string]bool, gates map[string]string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wires = map[string]bool{}
	gates = map[string]string{}

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	readWires := true
	for scanner.Scan() {
		line := scanner.Text()

		// we want to read the gates after first empty line
		if line == "" {
			readWires = false
			continue
		}

		if readWires {
			parts := strings.Split(line, ": ")
			wire := parts[0]
			if value, err := strconv.ParseBool(parts[1]); err != nil {
				log.Fatal("could not convert to bool", err)
			} else {
				wires[wire] = value
			}
		} else {
			parts := strings.Split(line, " -> ")
			instructions := parts[0]
			output := parts[1]
			gates[output] = instructions
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return wires, gates
}

func XOR(input1, input2 bool) bool {
	if input1 && input2 {
		return false
	}
	if !input1 && !input2 {
		return false
	}
	return true
}

func AND(input1, input2 bool) bool {
	if input1 && input2 {
		return true
	}
	return false
}

func OR(input1, input2 bool) bool {
	if input1 || input2 {
		return true
	}
	return false
}

var (
	memo = map[string]bool{}
)

func decodeInstructions(instruction string) (wire1 string, wire2 string, operation string) {
	parts := strings.Fields(instruction)
	wire1, operation, wire2 = parts[0], parts[1], parts[2]
	return wire1, wire2, operation
}

func evaluateWire(input string, gates map[string]string) (output bool) {
	// memoized already
	if output, exists := memo[input]; exists {
		return output
	}

	instruction := gates[input]

	// need to decode expression
	wire1, wire2, operation := decodeInstructions(instruction)
	input1 := evaluateWire(wire1, gates)
	input2 := evaluateWire(wire2, gates)

	switch operation {
	case "XOR":
		output = XOR(input1, input2)
	case "OR":
		output = OR(input1, input2)
	case "AND":
		output = AND(input1, input2)
	default:
		log.Fatal("could not find operation", operation)
	}
	memo[input] = output
	return output
}

func part1(gates map[string]string) int64 {
	zGates := []string{}
	for input := range gates {
		evaluateWire(input, gates)
		if strings.Contains(input, "z") {
			zGates = append(zGates, input)
		}
	}
	// will be sorted z00 first and higher later etc
	sort.Strings(zGates)

	binaryString := ""
	for _, zGate := range zGates {
		output := memo[zGate]
		// put the new bit on the left side of the binary
		if output {
			binaryString = "1" + binaryString
		} else {
			binaryString = "0" + binaryString
		}
	}

	decimalNumber, err := strconv.ParseInt(binaryString, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return decimalNumber
}

func main() {
	wires, gates := read_input()
	for wire, output := range wires {
		memo[wire] = output
	}
	decimalNumber := part1(gates)

	fmt.Println("part 1: ", decimalNumber)
	// fmt.Println("part 2: ")
}
