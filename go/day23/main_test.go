package main

import (
	"strings"
	"testing"
)

const smallExample = `.....
..##.
..#..
.....
..##.
.....`

const example = `....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..`

func Test_grove_stepN(t *testing.T) {
	grove := parseGrove(strings.ReplaceAll(example, "\n", "\r\n"))
	grove.stepN(10)
	res := grove.emptySpacesInRectangle(grove.smallestRectangle())
	if res != 110 {
		t.Errorf("Should return 110, not %d", res)
	}
}

func Test_grove_stepUntilNoChanges(t *testing.T) {
	grove := parseGrove(strings.ReplaceAll(example, "\n", "\r\n"))
	res := grove.stepUntilNoChanges()
	if res != 20 {
		t.Errorf("Should return 20, not %d", res)
	}
}
