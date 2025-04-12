package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

type clawMachine struct {
	buttonA coord
	buttonB coord
	prize   coord
}

func (this coord) move(other coord) coord {
	return coord{
		x: this.x + other.x,
		y: this.y + other.y,
	}
}

func readInput() (lines []string) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

func getClawMachines(lines []string) (machines []clawMachine) {
	for i := 0; i <= len(lines)-3; i += 4 {
		buttonA := parseSingleLine(lines[i])
		buttonB := parseSingleLine(lines[i+1])
		prize := parseSingleLine(lines[i+2])
		machine := clawMachine{buttonA, buttonB, prize}
		machines = append(machines, machine)
	}

	return machines
}

func parseSingleLine(line string) (xy coord) {
	coordString := strings.Split(line, ": ")[1]

	xy_string := strings.Split(coordString, ", ")
	x_string, y_string := xy_string[0], xy_string[1]
	// There will be X+{int} and Y+{int} respectively want
	// to remove the first to chars for this reason.
	x_string, y_string = x_string[2:], y_string[2:]

	x, err := strconv.Atoi(x_string)
	if err != nil {
		log.Fatal(err)
	}
	xy.x = x

	y, err := strconv.Atoi(y_string)
	if err != nil {
		log.Fatal(err)
	}
	xy.y = y

	return xy
}

func solveSecondOrderPolynomial(clawMachine clawMachine, max100Presses bool) (cost int) {
	// defined in the readme.
	cost_a := 3
	cost_b := 1

	a_x, b_x := clawMachine.buttonA.x, clawMachine.buttonB.x
	a_y, b_y := clawMachine.buttonA.y, clawMachine.buttonB.y
	p_x, p_y := clawMachine.prize.x, clawMachine.prize.y

	// We are dealing with the following equations
	// a* a_x + b * b_x = p_x
	// a* a_y + b * b_y = p_y
	// determinant != 0 will show if we can solve it or not.
	determinant := a_x*b_y - a_y*b_x

	if determinant == 0 {
		// If determinant is 0 then it will either be
		// insolvable or have infinite solutions given a < 0
		if a_x/a_y != b_x/b_y {
			return 0
		}
		if a_x/a_y != p_x/p_y {
			return 0
		}
		// Here we have infinite solutions, but since the
		// lowest a is 0 we will have the following solution
		b := p_x / b_x

		cost = b * cost_b
		return cost
	}

	// We need to have integer solutions, otherwise it is
	// not valid for the button presses. Golang will
	// automatically floor the result so we check if it is a
	// float this way.
	if (p_x*b_y-p_y*b_x)%determinant != 0 || (p_y*a_x-p_x*a_y)%determinant != 0 {
		return 0
	}

	a := (p_x*b_y - p_y*b_x) / determinant
	b := (p_y*a_x - p_x*a_y) / determinant

	// rule if any button is pressed 100 times then that is
	// counted as non-solveable used in part 1 but not in
	// part 2
	if max100Presses {
		if a > 100 || b > 100 {
			return 0
		}
	}

	cost = a*cost_a + b*cost_b

	return cost
}

func updatePrizeLocation(machines []clawMachine) (updatedMachines []clawMachine) {
	for _, machine := range machines {
		machine.prize.x += 10000000000000
		machine.prize.y += 10000000000000
		updatedMachines = append(updatedMachines, machine)
	}

	return updatedMachines
}

func main() {
	lines := readInput()
	machines := getClawMachines(lines)

	max100Presses := true
	total_cost := 0
	for _, machine := range machines {
		total_cost += solveSecondOrderPolynomial(machine, max100Presses)
	}
	fmt.Println("total cost part 1:", total_cost)

	// part 2 no longer constraint on the max number of
	// presses but updated prize location
	machines = updatePrizeLocation(machines)
	max100Presses = false
	total_cost = 0
	for _, machine := range machines {
		total_cost += solveSecondOrderPolynomial(machine, max100Presses)
	}
	fmt.Println("total cost part 2:", total_cost)
}
