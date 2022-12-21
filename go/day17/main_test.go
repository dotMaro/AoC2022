package main

import "testing"

const example = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func Test_chamber_highestRock(t *testing.T) {
	c := newChamber("<")
	res := c.highestRock()
	if res != 0 {
		t.Errorf("Should return 0, not %d", res)
	}
	c.grid = append(c.grid, 0b111)
	res = c.highestRock()
	if res != 1 {
		t.Errorf("Should return 1, not %d", res)
	}
}

func Test_chamber_addRocks(t *testing.T) {
	c := newChamber(example)
	res := c.addRocks(2022)
	if res != 3068 {
		t.Errorf("Should return 3068, not %d", res)
	}
}
