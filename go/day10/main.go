package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("go/day10/input.txt")
	cpu := parse(string(input))
	sum := cpu.cycleSum([]int{20, 60, 100, 140, 180, 220})
	fmt.Printf("Part 1. %d\n", sum)
	cpu.cycleUntilCompletion()
	fmt.Printf("Part 2.\n%s\n", cpu.String())
}

func parse(input string) cpu {
	var instructions []instruction
	for _, line := range strings.Split(input, "\r\n") {
		words := strings.SplitN(line, " ", 2)
		instructionType := parseInstructionType(words[0])
		value := 0
		if len(words) > 1 {
			value, _ = strconv.Atoi(words[1])
		}
		instruction := instruction{
			instructionType: instructionType,
			value:           value,
			cyclesRemaining: instructionType.cyclesRequired(),
		}
		instructions = append(instructions, instruction)
	}
	return cpu{
		x:                  1,
		currentCycle:       0,
		currentInstruction: nil,
		instructions:       instructions,
		pc:                 0,
		crt:                []bool{},
	}
}

type cpu struct {
	x                  int
	currentCycle       int
	currentInstruction *instruction
	instructions       []instruction
	pc                 int
	crt                []bool
}

func (c *cpu) cycleSum(cycles []int) int {
	sum := 0
	for _, cycle := range cycles {
		sum += c.cycleTo(cycle)
	}
	return sum
}

func (c *cpu) cycleTo(n int) int {
	var x int
	for c.currentCycle < n {
		x = c.cycle()
	}
	return x * c.currentCycle
}

func (c *cpu) cycleUntilCompletion() int {
	var x int
	for c.pc < len(c.instructions)-1 {
		x = c.cycle()
	}
	return x * c.currentCycle
}

func (c *cpu) cycle() int {
	if c.currentInstruction == nil {
		c.currentInstruction = &c.instructions[c.pc]
		c.pc++
	}

	c.crt = append(c.crt, c.currentCycle%40 >= c.x-1 && c.currentCycle%40 <= c.x+1)

	c.currentInstruction.cyclesRemaining--
	xBeforeExecute := c.x
	if c.currentInstruction.cyclesRemaining == 0 {
		c.executeInstruction(*c.currentInstruction)
		c.currentInstruction = nil
	}
	c.currentCycle++
	return xBeforeExecute
}

func (c *cpu) executeInstruction(instruction instruction) {
	switch instruction.instructionType {
	case noop:
		return
	case addx:
		c.x += instruction.value
	}
}

func (c *cpu) String() string {
	var b strings.Builder
	for i, crt := range c.crt {
		if i%40 == 0 {
			b.WriteRune('\n')
		}
		if crt {
			b.WriteRune('#')
		} else {
			b.WriteRune('.')
		}
	}
	return b.String()
}

type instruction struct {
	instructionType instructionType
	value           int
	cyclesRemaining int
}

type instructionType int

const (
	noop instructionType = iota
	addx
)

func parseInstructionType(s string) instructionType {
	switch s {
	case "noop":
		return noop
	case "addx":
		return addx
	default:
		panic(fmt.Sprintf("invalid instruction type %q", s))
	}
}

func (t instructionType) cyclesRequired() int {
	switch t {
	case noop:
		return 1
	case addx:
		return 2
	default:
		panic("invalid instruction type")
	}
}
