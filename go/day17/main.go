package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day17/input.txt")
	if err != nil {
		panic(err)
	}
	chamber := newChamber(string(input))
	highest := chamber.addRocks(2022)
	fmt.Printf("Part 1. After 2022 fallen rocks the top rock is at %d\n", highest)
	chamber = newChamber(string(input))
	highest = chamber.addRocks(1_000_000_000_000)
	fmt.Printf("Part 2. After 1000000000000 fallen rocks the top rock is at %d\n", highest)
}

type shape []uint8

var shapes = []shape{
	// ####
	{0b00_1111_00},
	// .#.
	// ###
	// .#.
	{
		0b00_010_000,
		0b00_111_000,
		0b00_010_000,
	},
	// ..#
	// ..#
	// ###
	{
		0b00_001_000,
		0b00_001_000,
		0b00_111_000,
	},
	// #
	// #
	// #
	// #
	{
		0b00_1_00000,
		0b00_1_00000,
		0b00_1_00000,
		0b00_1_00000,
	},
	// ##
	// ##
	{
		0b00_11_0000,
		0b00_11_0000,
	},
}

type chamber struct {
	grid           []uint8
	gridOffset     uint64
	fallingShape   shape
	fallingRockY   int
	rockCount      int
	directions     []direction
	directionIndex int
	shapeIndex     int
	history        map[signature]result
	startIndexes   signature
}

type coords struct {
	x, y int
}

type direction int

const (
	left direction = iota
	right
)

type signature struct {
	top            string
	directionIndex int
	shapeIndex     int
}

type result struct {
	height int
	rocks  uint64
}

func newChamber(input string) chamber {
	var directions []direction
	for _, c := range input {
		var d direction
		switch c {
		case '<':
			d = left
		case '>':
			d = right
		}
		directions = append(directions, d)
	}
	return chamber{
		grid:           []uint8{0b1111111},
		gridOffset:     0,
		fallingShape:   nil,
		fallingRockY:   0,
		rockCount:      0,
		directions:     directions,
		directionIndex: 0,
		shapeIndex:     0,
		history:        make(map[signature]result),
		startIndexes: signature{
			directionIndex: 0,
			shapeIndex:     0,
		},
	}
}

func (c *chamber) expandTo(n int) {
	for len(c.grid) < n+1 {
		c.grid = append(c.grid, 0)
	}
}

func (c *chamber) topAsString() string {
	var top strings.Builder
	for _, row := range c.grid[max(len(c.grid)-30, 0):] {
		top.WriteString(fmt.Sprint(row))
		top.WriteRune(',')
	}
	return top.String()
}

func (c *chamber) recordSignature(rockCount uint64) {

	sign := signature{
		top:            c.topAsString(),
		directionIndex: c.directionIndex,
		shapeIndex:     c.shapeIndex,
	}
	c.history[sign] = result{
		height: c.highestRock(),
		rocks:  rockCount,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *chamber) addRocks(n uint64) uint64 {
	var rock uint64
	for rock = 0; rock < n; rock++ {
		shapeCopy := make(shape, len(shapes[c.shapeIndex]))
		copy(shapeCopy, shapes[c.shapeIndex])
		c.fallingShape = shapeCopy
		c.shapeIndex = (c.shapeIndex + 1) % len(shapes)
		y := c.highestRock() + 4 + len(c.fallingShape) - 1
		c.fallingRockY = y

		falling := true
		for falling {
			direction := c.directions[c.directionIndex]
			c.directionIndex = (c.directionIndex + 1) % len(c.directions)
			if c.canGo(direction) {
				switch direction {
				case left:
					for i, row := range c.fallingShape {
						c.fallingShape[i] = row << 1
					}
				case right:
					for i, row := range c.fallingShape {
						c.fallingShape[i] = row >> 1
					}
				}
			}

			if c.canGoDown() {
				c.fallingRockY--
			} else {
				c.insertFallingShapeIntoGrid(n, rock)
				falling = false
			}
		}

		if rock > 1000 {
			sign := signature{
				top:            c.topAsString(),
				directionIndex: c.directionIndex,
				shapeIndex:     c.shapeIndex,
			}
			res, alreadySeen := c.history[sign]
			if alreadySeen {
				rockDiff := rock - res.rocks
				heightDiff := c.highestRock() - res.height
				// The amount of times we can fit the recorded history until we reach the goal.
				amount := (n - rock) / rockDiff
				c.gridOffset += amount * uint64(heightDiff)
				rock += amount * rockDiff
			}
		}
		c.recordSignature(rock)
	}

	return uint64(c.highestRock()) + c.gridOffset
}

func (c *chamber) insertFallingShapeIntoGrid(n, i uint64) {
	c.expandTo(c.fallingRockY)

	for localY, row := range c.fallingShape {
		c.grid[c.fallingRockY-localY] |= row
	}

	// c.trimIfPossible(n, i)
}

// func (c *chamber) trimIfPossible(n, i uint64) {
// 	for localY, row := range c.grid[c.fallingRockY-len(*c.fallingShape)+1 : c.fallingRockY] {
// 		if row&0b1111111_0 != 0b1111111_0 {
// 			continue
// 		}

// 		globalY := localY + c.fallingRockY - len(*c.fallingShape) + 1
// 		// fmt.Printf("Full row at %d. direction: %v. shape: %v\n", globalY, c.directionIndex, c.shapeIndex)
// 		highestRow := localY == 0
// 		indexes := signature{
// 			directionIndex: c.directionIndex,
// 			shapeIndex:     c.shapeIndex,
// 		}
// 		// Register what score we got with the old startIndexes. If we ever get the same startIndexes then we'll
// 		// already know what score and index offsets they'll get.
// 		// Only do it if we're at the highest row though, so there are identical prerequisites.
// 		if highestRow {
// 			res, alreadyRegistered := c.history[indexes]
// 			if alreadyRegistered && uint64(res.rocks)+i < n {
// 				// c.directionIndex = res.indexes.directionIndex
// 				// c.shapeIndex = res.indexes.shapeIndex

// 				cont, hasCont := c.history[res.indexes]
// 				if hasCont {
// 					c.history[indexes] = cont
// 				}

// 				fmt.Printf("History match! %d %#v %#v\n", globalY, indexes, res)
// 			} else {
// 				c.history[c.startIndexes] = result{
// 					height:  globalY,
// 					rocks:   c.rockCount,
// 					indexes: indexes,
// 				}
// 			}
// 		}
// 		c.startIndexes = indexes
// 		c.rockCount = 0

// 		c.grid = c.grid[globalY:]
// 		c.gridOffset += uint64(globalY)
// 		break
// 	}
// }

func (c *chamber) highestRock() int {
	for y := len(c.grid) - 1; y >= 0; y-- {
		row := c.grid[y]
		if row != 0 {
			return y
		}
	}
	return 0
}

func (c *chamber) canGo(d direction) bool {
	for localY, row := range c.fallingShape {
		switch d {
		case left:
			if row&0b1000000_0 != 0 || c.fallingRockY-localY < len(c.grid) && row<<1&c.grid[c.fallingRockY-localY] != 0 {
				return false
			}
		case right:
			if row&0b0000001_0 != 0 || c.fallingRockY-localY < len(c.grid) && row>>1&c.grid[c.fallingRockY-localY] != 0 {
				return false
			}
		}
	}
	return true
}

func (c *chamber) canGoDown() bool {
	for localY, row := range c.fallingShape {
		if c.fallingRockY-localY > len(c.grid) {
			continue
		}
		if row&c.grid[c.fallingRockY-localY-1] != 0 {
			return false
		}
	}
	return true
}

func (c *chamber) String() string {
	var b strings.Builder
	for y := len(c.grid) - 1; y > 0; y-- {
		row := c.grid[y]
		b.WriteRune('|')
		for i := 6; i >= 0; i-- {
			if row&(0b0000001_0<<i) != 0 {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteString("|\n")
	}
	b.WriteString("+-------+")
	return b.String()
}
