package main

import (
	"strings"
	"testing"
)

const example = `#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`

func Test_valley_findShortestPath(t *testing.T) {
	valley := parseValley(strings.ReplaceAll(example, "\n", "\r\n"))
	res := valley.findShortestPath(false)
	if res != 18 {
		t.Errorf("Should return 18, not %d", res)
	}
	res = valley.findShortestPath(true)
	if res != 54 {
		t.Errorf("Should return 54, not %d", res)
	}
}
