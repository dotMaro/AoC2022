package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day25/input.txt")
	if err != nil {
		panic(err)
	}
	sum := sum(string(input))
	snafu := intToSnafu(sum)
	fmt.Printf("The SNAFU of the sum is %s", snafu)
}

func sum(input string) int {
	sum := 0
	for _, line := range strings.Split(input, "\r\n") {
		sum += snafuToInt(line)
	}
	return sum
}

func snafuToInt(snafu string) int {
	total := 0
	p := 1
	for i := len(snafu) - 1; i >= 0; i-- {
		switch snafu[i] {
		case '2':
			total += 2 * p
		case '1':
			total += p
		case '0': // Noop.
		case '-':
			total -= p
		case '=':
			total -= 2 * p
		}
		p *= 5
	}
	return total
}

func intToSnafu(i int) string {
	var snafu strings.Builder
	total := 0
	p := 1
	digits := 0
	for i > total {
		total += 2 * p
		p *= 5
		digits++
	}
	for digits > 0 {
		pp := p / 5
		if total-4*pp >= i {
			total -= 4 * pp
			snafu.WriteRune('=')
		} else if total-3*pp >= i {
			total -= 3 * pp
			snafu.WriteRune('-')
		} else if total-2*pp >= i {
			total -= 2 * pp
			snafu.WriteRune('0')
		} else if total-pp >= i {
			total -= pp
			snafu.WriteRune('1')
		} else {
			snafu.WriteRune('2')
		}
		p /= 5
		digits--
	}
	return snafu.String()
}
