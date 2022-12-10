package main

import (
	"strings"
	"testing"
)

const example = `addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop`

func Test_cpu_cycleTo(t *testing.T) {
	cpu := parse(strings.ReplaceAll(example, "\n", "\r\n"))
	testCases := []struct {
		cycles int
		exp    int
	}{
		{20, 420},
		{60, 1140},
		{cycles: 100, exp: 1800},
		{140, 2940},
		{180, 2880},
		{220, 3960},
	}
	for _, tc := range testCases {
		res := cpu.cycleTo(tc.cycles)
		if res != tc.exp {
			t.Errorf("Should return %d after cycle %d, but returned %d", tc.exp, tc.cycles, res)
		}
	}
}

func Test_cpu_cycleSum(t *testing.T) {
	cpu := parse(strings.ReplaceAll(example, "\n", "\r\n"))
	sum := cpu.cycleSum([]int{20, 60, 100, 140, 180, 220})
	exp := 13140
	if sum != exp {
		t.Errorf("Sum should be %d, not %d", exp, sum)
	}
}
