package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day21/input.txt")
	if err != nil {
		panic(err)
	}
	monkeys := parseMonkeys(string(input))
	fmt.Printf("Part 1. The root monkey has value %d\n", monkeys["root"].value(monkeys, 0))
	monkeys["root"] = operator{
		monkey1Name: monkeys["root"].(operator).monkey1Name,
		monkey2Name: monkeys["root"].(operator).monkey2Name,
		operation:   equals,
	}
	monkeys["humn"] = wildcard{}
	monkeys["root"].value(monkeys, 0) // This will set the humn monkey to a value.
	fmt.Printf("Part 2. The humn monkey should have value %d\n", monkeys["humn"])
}

func parseMonkeys(input string) monkeys {
	monkeys := make(monkeys)
	for _, line := range strings.Split(input, "\r\n") {
		words := strings.Split(line, " ")
		name := words[0][:len(words[0])-1]
		if len(words) == 2 {
			val, err := strconv.Atoi(words[1])
			if err != nil {
				panic(err)
			}
			monkeys[name] = number(val)
		} else {
			var operation operation
			switch words[2][0] {
			case '+':
				operation = add
			case '-':
				operation = subtract
			case '*':
				operation = multiply
			case '/':
				operation = divide
			}
			monkeys[name] = operator{
				monkey1Name: words[1],
				monkey2Name: words[3],
				operation:   operation,
			}
		}
	}
	return monkeys
}

type monkeys map[string]value

type value interface {
	value(map[string]value, int) int
	containsWildcard(map[string]value) bool
}

type number int

func (n number) value(monkeys map[string]value, shouldBeEqualTo int) int {
	return int(n)
}

func (n number) containsWildcard(map[string]value) bool {
	return false
}

type wildcard struct{}

func (w wildcard) value(monkeys map[string]value, shouldBeEqualTo int) int {
	monkeys["humn"] = number(shouldBeEqualTo)
	return shouldBeEqualTo
}

func (w wildcard) containsWildcard(map[string]value) bool {
	return true
}

type operator struct {
	monkey1Name, monkey2Name string
	operation
}

func (o operator) value(monkeys map[string]value, shouldBeEqualTo int) int {
	return o.operation(monkeys, o.monkey1Name, o.monkey2Name, shouldBeEqualTo)
}

func (o operator) containsWildcard(monkeys map[string]value) bool {
	return monkeys[o.monkey1Name].containsWildcard(monkeys) || monkeys[o.monkey2Name].containsWildcard(monkeys)
}

type operation func(monkeys map[string]value, operand1, operand2 string, shouldBeEqualTo int) int

func add(monkeys map[string]value, monkey1Name, monkey2Name string, shouldBeEqualTo int) int {
	monkey1 := monkeys[monkey1Name]
	monkey2 := monkeys[monkey2Name]
	if monkey1.containsWildcard(monkeys) {
		monkey2Value := monkey2.value(monkeys, 0)
		return monkey1.value(monkeys, shouldBeEqualTo-monkey2Value) + monkey2Value
	} else if monkey2.containsWildcard(monkeys) {
		monkey1Value := monkey1.value(monkeys, 0)
		return monkey1Value + monkey2.value(monkeys, shouldBeEqualTo-monkey1Value)
	}

	return monkey1.value(monkeys, 0) + monkey2.value(monkeys, 0)
}

func subtract(monkeys map[string]value, monkey1Name, monkey2Name string, shouldBeEqualTo int) int {
	monkey1 := monkeys[monkey1Name]
	monkey2 := monkeys[monkey2Name]
	if monkey1.containsWildcard(monkeys) {
		monkey2Value := monkey2.value(monkeys, 0)
		return monkey1.value(monkeys, monkey2Value+shouldBeEqualTo) - monkey2Value
	} else if monkey2.containsWildcard(monkeys) {
		monkey1Value := monkey1.value(monkeys, 0)
		return monkey1Value - monkey2.value(monkeys, monkey1Value-shouldBeEqualTo)
	}

	return monkey1.value(monkeys, 0) - monkey2.value(monkeys, 0)
}

func multiply(monkeys map[string]value, monkey1Name, monkey2Name string, shouldBeEqualTo int) int {
	monkey1 := monkeys[monkey1Name]
	monkey2 := monkeys[monkey2Name]
	if monkey1.containsWildcard(monkeys) {
		monkey2Value := monkey2.value(monkeys, 0)
		return monkey1.value(monkeys, shouldBeEqualTo/monkey2Value) / monkey2Value
	} else if monkey2.containsWildcard(monkeys) {
		monkey1Value := monkey1.value(monkeys, 0)
		return monkey1Value / monkey2.value(monkeys, shouldBeEqualTo/monkey1Value)
	}

	return monkey1.value(monkeys, 0) * monkey2.value(monkeys, 0)
}

func divide(monkeys map[string]value, monkey1Name, monkey2Name string, shouldBeEqualTo int) int {
	monkey1 := monkeys[monkey1Name]
	monkey2 := monkeys[monkey2Name]
	if monkey1.containsWildcard(monkeys) {
		monkey2Value := monkey2.value(monkeys, 0)
		return monkey1.value(monkeys, monkey2Value*shouldBeEqualTo) / monkey2Value
	} else if monkey2.containsWildcard(monkeys) {
		monkey1Value := monkey1.value(monkeys, 0)
		return monkey1Value / monkey2.value(monkeys, monkey1Value/shouldBeEqualTo)
	}

	return monkey1.value(monkeys, 0) / monkey2.value(monkeys, 0)
}

func equals(monkeys map[string]value, monkey1Name, monkey2Name string, shouldBeEqualTo int) int {
	monkey1 := monkeys[monkey1Name]
	monkey2 := monkeys[monkey2Name]
	if monkey1.containsWildcard(monkeys) {
		monkey2Value := monkey2.value(monkeys, 0)
		return monkey1.value(monkeys, monkey2Value)
	}
	monkey1Value := monkey1.value(monkeys, 0)
	return monkey2.value(monkeys, monkey1Value)
}
