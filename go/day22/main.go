package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day22/input.txt")
	if err != nil {
		panic(err)
	}
	board := parseBoard(string(input))
	board.followInstructions(false)
	fmt.Printf("Part 1. The password is %d\n", board.password())
	board = parseBoard(string(input))
	board.followInstructions(true)
	fmt.Printf("Part 2. The password is %d\n", board.password())
}

func parseBoard(input string) board {
	grid := make([][]tile, 0, 200)
	lines := strings.Split(input, "\r\n")
	for _, line := range lines {
		if line == "" {
			break
		}
		row := make([]tile, len(lines[0]))
		for i, c := range line {
			row[i] = parseTile(c)
		}
		grid = append(grid, row)
	}

	var curNumber string
	instructions := make([]instruction, 0, 100)
	for _, c := range lines[len(lines)-1] {
		if c >= '0' && c <= '9' {
			curNumber += string(c)
			continue
		}
		steps, err := strconv.Atoi(curNumber)
		if err != nil {
			panic(err)
		}
		curNumber = ""
		stepInstruction := instruction{
			steps: &steps,
		}
		dir := parseDirection(c)
		directionInstruction := instruction{
			direction: &dir,
		}
		instructions = append(instructions, stepInstruction, directionInstruction)
	}
	if curNumber != "" {
		steps, err := strconv.Atoi(curNumber)
		if err != nil {
			panic(err)
		}
		stepInstruction := instruction{
			steps: &steps,
		}
		instructions = append(instructions, stepInstruction)
	}
	var xPos int
	for x, tile := range grid[0] {
		if tile == open {
			xPos = x
			break
		}
	}
	return board{
		grid:         grid,
		instructions: instructions,
		curPos:       coord{x: xPos, y: 0},
		curDir:       east,
	}
}

type board struct {
	grid         [][]tile
	instructions []instruction
	curPos       coord
	curDir       cardinalDirection
}

type instruction struct {
	steps     *int
	direction *direction
}

type coord struct {
	x, y int
}

func parseDirection(r rune) direction {
	var dir direction
	switch r {
	case 'R':
		dir = right
	case 'L':
		dir = left
	default:
		panic("invalid direction rune")
	}
	return dir
}

type direction uint8

const (
	right direction = iota
	left
)

type cardinalDirection uint8

const (
	east cardinalDirection = iota
	south
	west
	north
)

func parseTile(c rune) tile {
	var tile tile
	switch c {
	case ' ':
		tile = void
	case '.':
		tile = open
	case '#':
		tile = wall
	default:
		panic("invalid tile char")
	}
	return tile
}

type tile uint8

const (
	void tile = iota
	open
	wall
)

func (b *board) tile(c coord) tile {
	return b.grid[c.y][c.x]
}

func (b *board) followInstructions(cubeMode bool) {
	for _, instruction := range b.instructions {
		if instruction.steps != nil {
			b.stepForward(*instruction.steps, cubeMode)
			continue
		}
		b.adjustDirection(*instruction.direction)
	}
}

func (b *board) stepForward(n int, cubeMode bool) {
	for i := 0; i < n; i++ {
		nextPos, dir := b.forwardFunc(cubeMode)(b.curPos, b.curDir)
		for b.tile(nextPos) == void {
			nextPos, dir = b.forwardFunc(cubeMode)(nextPos, dir)
		}
		if b.tile(nextPos) == wall {
			return
		}
		b.curPos = nextPos
		b.curDir = dir
	}
}

func (b *board) forwardFunc(cubeMode bool) func(coord, cardinalDirection) (coord, cardinalDirection) {
	borderCheck := func(c coord, d cardinalDirection) (coord, cardinalDirection) {
		// Top right.
		if c.x == 150 && c.y < 50 {
			coord := coord{
				x: 99,
				y: 149 - c.y,
			}
			return coord, west
		} else if c.x == 49 && c.y < 100 {
			// Above center left.
			if c.y >= 50 {
				coord := coord{
					x: c.y - 50,
					y: 100,
				}
				return coord, south
			}
			// Top left.
			if c.y < 50 {
				coord := coord{
					x: 0,
					y: 149 - c.y,
				}
				return coord, east
			}
		} else if c.x == -1 {
			// Middle left.
			if c.y >= 100 && c.y < 150 {
				coord := coord{
					x: 50,
					y: 49 - (c.y - 100),
				}
				return coord, east
			}
			// Bottom left.
			if c.y >= 150 {
				coord := coord{
					x: c.y - 100,
					y: 0,
				}
				return coord, south
			}
		} else if c.x == 100 {
			// Above center right.
			if d == east && c.y >= 50 && c.y < 100 {
				coord := coord{
					x: c.y + 50,
					y: 49,
				}
				return coord, north
			}
			if c.y >= 100 && c.y < 150 {
				// Middle right.
				coord := coord{
					x: 149,
					y: 49 - (c.y - 100),
				}
				return coord, west
			}
		} else if d == east && c.x == 50 && c.y >= 150 {
			// Bottom right.
			coord := coord{
				x: c.y - 100,
				y: 149,
			}
			return coord, north
		} else if c.y == -1 {
			if c.x < 100 {
				// Top up.
				coord := coord{
					x: 0,
					y: 100 + c.x,
				}
				return coord, east
			}
			// Right top up.
			coord := coord{
				x: c.x - 100,
				y: 199,
			}
			return coord, north
		} else if d == north && c.y == 99 && c.x < 50 {
			// Left up.
			coord := coord{
				x: 50,
				y: 50 + c.x,
			}
			return coord, east
		} else if d == south && c.y == 50 && c.x >= 100 {
			// Top right down.
			coord := coord{
				x: 99,
				y: c.x - 50,
			}
			return coord, west
		} else if d == south && c.y == 150 && c.x >= 50 {
			// Center down.
			coord := coord{
				x: 49,
				y: c.x + 100,
			}
			return coord, west
		} else if c.y == 200 {
			// Bottom down.
			coord := coord{
				x: c.x + 100,
				y: 0,
			}
			return coord, south
		}

		return c, d
	}

	var forwardFunc func(coord, cardinalDirection) (coord, cardinalDirection)
	switch b.curDir {
	case north:
		if cubeMode {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				return borderCheck(coord{x: c.x, y: c.y - 1}, d)
			}
		} else {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				if c.y == 0 {
					c.y = len(b.grid)
				}
				return coord{x: c.x, y: c.y - 1}, d
			}
		}
	case east:
		if cubeMode {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				return borderCheck(coord{x: c.x + 1, y: c.y}, d)
			}
		} else {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				return coord{x: (c.x + 1) % len(b.grid[0]), y: c.y}, d
			}
		}
	case south:
		if cubeMode {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				// return borderCheck(coord{x: c.x, y: (c.y + 1) % len(b.grid)}, d)
				return borderCheck(coord{x: c.x, y: c.y + 1}, d)
			}
		} else {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				return coord{x: c.x, y: (c.y + 1) % len(b.grid)}, d
			}

		}
	case west:
		if cubeMode {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				return borderCheck(coord{x: c.x - 1, y: c.y}, d)
			}
		} else {
			forwardFunc = func(c coord, d cardinalDirection) (coord, cardinalDirection) {
				if c.x == 0 {
					c.x = len(b.grid[0])
				}
				return coord{x: c.x - 1, y: c.y}, d
			}

		}
	}
	return forwardFunc
}

func (b *board) adjustDirection(d direction) {
	var newDir cardinalDirection
	if d == right {
		newDir = cardinalDirection((uint8(b.curDir) + 1) % 4)
	} else {
		if uint8(b.curDir) == 0 {
			newDir = cardinalDirection(3)
		} else {
			newDir = cardinalDirection(int8(b.curDir) - 1)
		}
	}
	b.curDir = newDir
}

func (b *board) password() int {
	return 1000*(b.curPos.y+1) + 4*(b.curPos.x+1) + int(b.curDir)
}
