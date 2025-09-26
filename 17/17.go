package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() (int, int, int, []int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var instructions []int
	var regA, regB, regC int

	// Scan and format lines.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check registers
		if strings.Contains(line, "A:") {
			regA, err = strconv.Atoi(strings.Split(line, "A: ")[1])
		}
		if strings.Contains(line, "B:") {
			regB, err = strconv.Atoi(strings.Split(line, "B: ")[1])
		}
		if strings.Contains(line, "C:") {
			regC, err = strconv.Atoi(strings.Split(line, "C: ")[1])
		}
		if !strings.Contains(line, "Program:") {
			continue
		}
		// Return program
		cleanedString := strings.Split(line, "Program: ")[1]
		stringInstructions := strings.Split(cleanedString, ",")
		for _, stringInstruction := range stringInstructions {
			instruction, _ := strconv.Atoi(stringInstruction)
			instructions = append(instructions, instruction)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return regA, regB, regC, instructions
}

func modulo(x int, d int) int {
	// Unfortunately Golang's % operator is not modulo
	// at least not in the common sense but instead
	// remainder.
	m := x % d
	if x < 0 && d < 0 {
		m -= d
	}
	if x < 0 && d > 0 {
		m += d
	}
	return m
}

type computer struct {
	A            int
	B            int
	C            int
	instructions []int
	pointer      int
	output       []int
}

func (current computer) swapAregister(newA int) (next computer) {
	return computer{
		A:            newA,
		B:            current.B,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer,
		output:       current.output,
	}
}

func (this computer) literalOperandToComboOperand(operand int) (input int) {
	if (operand < 0) || (6 < operand) {
		log.Fatal("Combo operand is something else than 0-7")
	}
	if operand < 4 {
		return operand
	}
	if operand == 4 {
		return this.A
	}
	if operand == 5 {
		return this.B
	}
	if operand == 6 {
		return this.C
	}
	return
}

func (current computer) op0(operand int) (next computer) {
	// A is divided by 2 to the power of combo operand and stored in A register
	input := current.literalOperandToComboOperand(operand)
	// in go bitwise shift can be used to calculate y/2^X by doin (y >> x)
	output := current.A >> input
	return computer{
		A:            output,
		B:            current.B,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func (current computer) op1(operand int) (next computer) {
	// B is the bitwise xor with itself and the input
	input := operand
	output := current.B ^ input
	return computer{
		A:            current.A,
		B:            output,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func (current computer) op2(operand int) (next computer) {
	// Combo operand modulo 8 is written to B
	input := current.literalOperandToComboOperand(operand)
	output := modulo(input, 8)
	return computer{
		A:            current.A,
		B:            output,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func (current computer) op3(operand int) (next computer) {
	// set pointer to literal operand unless A register = 0
	input := operand
	var output int
	if current.A == 0 {
		// does nothing if A register = 0
		output = current.pointer + 2
	} else {
		// jumps by literal operand if A register != 0
		output = input
	}
	return computer{
		A:            current.A,
		B:            current.B,
		C:            current.C,
		instructions: current.instructions,
		pointer:      output,
		output:       current.output,
	}
}

func (current computer) op4(operand int) (next computer) {
	// B is the bitwise xor of B and C operand is not used
	output := current.B ^ current.C
	return computer{
		A:            current.A,
		B:            output,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func (current computer) op5(operand int) (next computer) {
	// modulo 8 of combo operator added to outputs
	input := current.literalOperandToComboOperand(operand)
	output := modulo(input, 8)
	return computer{
		A:            current.A,
		B:            current.B,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       append(current.output, output),
	}
}

func (current computer) op6(operand int) (next computer) {
	// A is divided by combo operand to the power of 2 and stored in B register
	input := current.literalOperandToComboOperand(operand)
	// in go bitwise shift can be used to calculate y/2^X by doin (y >> x)
	output := current.A >> input
	return computer{
		A:            current.A,
		B:            output,
		C:            current.C,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func (current computer) op7(operand int) (next computer) {
	// A is divided by combo operand to the power of 2 and stored in C register
	input := current.literalOperandToComboOperand(operand)
	// in go bitwise shift can be used to calculate y/2^X by doin (y >> x)
	output := current.A >> input
	return computer{
		A:            current.A,
		B:            current.B,
		C:            output,
		instructions: current.instructions,
		pointer:      current.pointer + 2,
		output:       current.output,
	}
}

func completeNextInstruction(current computer) (next computer, completed bool) {
	if current.pointer >= len(current.instructions)-1 {
		return current, true
	}

	operation := current.instructions[current.pointer]
	operand := current.instructions[current.pointer+1]

	switch operation {
	case 0:
		next = current.op0(operand)
	case 1:
		next = current.op1(operand)
	case 2:
		next = current.op2(operand)
	case 3:
		next = current.op3(operand)
	case 4:
		next = current.op4(operand)
	case 5:
		next = current.op5(operand)
	case 6:
		next = current.op6(operand)
	case 7:
		next = current.op7(operand)
	default:
		log.Fatal("operation not in set 0-7")
	}
	return next, false
}

func intArrayToString(input []int) string {
	var stringArray []string
	for _, elem := range input {
		stringArray = append(stringArray, strconv.Itoa(elem))
	}
	return strings.Join(stringArray, ",")
}

func getOutputFromComputer(currentComputer computer) []int {
	// complete instructions
	instructionsFinished := false

	for !instructionsFinished {
		currentComputer, instructionsFinished = completeNextInstruction(currentComputer)
	}
	return currentComputer.output
}

func isMatchingSlice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func backtrackOutputWhichReflectsComputer(currentComputer computer) int {
	// the number to guess, will be large as a will be divided by 8 for each round and we need 16 outputs
	nextA := 0
	var nextComputer computer
	var nextOutput []int

	// start with finding how the lowest number should be backtracked and go from there
	for i := range len(currentComputer.instructions) {
		revI := len(currentComputer.instructions) - i - 1
		nextA = nextA << 3 // shift to left tree times => * 8 => ready to guess the next number
		nextComputer = currentComputer.swapAregister(nextA)
		nextOutput = getOutputFromComputer(nextComputer)

		// increase until we are ready for the next number
		for !isMatchingSlice(nextOutput, currentComputer.instructions[revI:]) {
			nextA++
			nextComputer = currentComputer.swapAregister(nextA)
			nextOutput = getOutputFromComputer(nextComputer)
		}

	}
	return nextA
}

func main() {
	regA, regB, regC, instructions := read_input()
	// init computer
	currentComputer := computer{
		A:            regA,
		B:            regB,
		C:            regC,
		instructions: instructions,
		pointer:      0,
		output:       []int{},
	}
	// format output and print
	p1_output := intArrayToString(getOutputFromComputer(currentComputer))
	fmt.Println("part 1: ", p1_output)

	p2_output := backtrackOutputWhichReflectsComputer(currentComputer)
	fmt.Println("part 2: ", p2_output)
}
