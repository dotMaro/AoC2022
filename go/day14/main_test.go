package main

import (
	"testing"
)

const example = "498,4 -> 498,6 -> 496,6\r\n503,4 -> 502,4 -> 502,9 -> 494,9"

func Test_cavern_addSandUntilOneLandsInTheAbyss(t *testing.T) {
	cavern := parseCavern(example, false)
	res := cavern.addSandUntilOneLandsInTheAbyss()
	if res != 24 {
		t.Errorf("Should return 24, but returned %d", res)
	}
}

func Test_cavern_addSandUntilSourceIsCovered(t *testing.T) {
	cavern := parseCavern(example, true)
	res := cavern.addSandUntilSourceIsCovered()
	if res != 93 {
		t.Errorf("Should return 93, but returned %d", res)
	}
}
