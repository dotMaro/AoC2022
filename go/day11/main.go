package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("go/day11/input.txt")

	monkeys := parseMonkeys(string(input))
	monkeyBusiness := rounds(monkeys, 20, true)
	fmt.Printf("Part 1. Monkey business is %d after 20 rounds.\n", monkeyBusiness)
	monkeys = parseMonkeys(string(input))
	monkeyBusiness = rounds(monkeys, 10000, false)
	fmt.Printf("Part 2. Monkey business is %d after 10000 rounds.\n", monkeyBusiness)
}

func parseMonkeys(input string) []*monkey {
	var monkeys []*monkey
	// Split by double newline, which is only used between monkey descriptions.
	for _, group := range strings.Split(input, "\r\n\r\n") {
		lines := strings.Split(group, "\r\n")
		var items []int
		for _, item := range strings.Split(lines[1][len("  Starting items: "):], ", ") {
			a, _ := strconv.Atoi(item)
			items = append(items, a)
		}
		operandChar := lines[2][len("  Operation: new = old ")]
		modifier := lines[2][len("  Operation: new = old ")+2:]
		operation := parseOperation(operandChar, modifier)
		testDivisible, _ := strconv.Atoi(lines[3][len("  Test: divisible by "):])
		testTrueTarget, _ := strconv.Atoi(lines[4][len("    If true: throw to monkey "):])
		testFalseTarget, _ := strconv.Atoi(lines[5][len("    If false: throw to monkey "):])
		monkey := monkey{
			items:           items,
			operation:       operation,
			testDivisible:   testDivisible,
			testTrueTarget:  testTrueTarget,
			testFalseTarget: testFalseTarget,
			inspectCount:    0,
		}
		monkeys = append(monkeys, &monkey)
	}
	return monkeys
}

func parseOperation(operandChar byte, modifier string) operation {
	var operation operation
	switch operandChar {
	case '+':
		if modifier == "old" {
			operation = func(old int) int {
				return old + old
			}
		} else {
			term, _ := strconv.Atoi(modifier)
			operation = func(old int) int {
				return old + term
			}
		}
	case '*':
		if modifier == "old" {
			operation = func(old int) int {
				return old * old
			}
		} else {
			factor, _ := strconv.Atoi(modifier)
			operation = func(old int) int {
				return old * factor
			}
		}
	}

	return operation
}

type monkey struct {
	items           []int
	operation       operation
	testDivisible   int
	testTrueTarget  int
	testFalseTarget int
	inspectCount    int
}

type operation func(old int) int

func (m *monkey) turn(monkeys []*monkey, productOfAllTests int, divide bool) {
	for _, item := range m.items {
		// Inspect.
		m.inspectCount++
		item = m.operation(item)
		if divide {
			item /= 3
		}
		// Test.
		item %= productOfAllTests
		if item%m.testDivisible == 0 {
			monkeys[m.testTrueTarget].items = append(monkeys[m.testTrueTarget].items, item)
		} else {
			monkeys[m.testFalseTarget].items = append(monkeys[m.testFalseTarget].items, item)
		}
	}
	// All items have been thrown.
	m.items = []int{}
}

func rounds(monkeys []*monkey, rounds int, divide bool) int {
	productOfAllTests := 1
	for _, m := range monkeys {
		productOfAllTests *= m.testDivisible
	}

	for i := 0; i < rounds; i++ {
		for _, monkey := range monkeys {
			monkey.turn(monkeys, productOfAllTests, divide)
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspectCount > monkeys[j].inspectCount
	})
	return monkeys[0].inspectCount * monkeys[1].inspectCount
}
