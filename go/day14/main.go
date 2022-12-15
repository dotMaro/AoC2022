package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("go/day14/input.txt")
	cavern := parseCavern(string(input), false)
	count := cavern.addSandUntilOneLandsInTheAbyss()
	fmt.Printf("Part 1. %d units of sand came to rest until one fell into the abyss\n", count)
	cavern = parseCavern(string(input), true)
	count = cavern.addSandUntilSourceIsCovered()
	fmt.Printf("Part 2. %d units of sand came to rest until the source got blocked\n", count)
}

func parseCavern(input string, withFloor bool) cavern {
	cavern := make(cavern, 200)
	for x := range cavern {
		cavern[x] = make([]tile, 1000)
	}

	highestY := 0
	for _, line := range strings.Split(input, "\r\n") {
		coords := strings.Split(line, " -> ")
		curCoord := parseCoord(coords[0])
		for _, coordStr := range coords[1:] {
			targetCoord := parseCoord(coordStr)
			if targetCoord.y > highestY {
				highestY = targetCoord.y
			}
			changed := true
			for changed {
				cavern[curCoord.y][curCoord.x] = rock
				curCoord, changed = curCoord.toward(targetCoord)
			}
		}
	}

	if withFloor {
		y := highestY + 2
		for x := range cavern[y] {
			cavern[y][x] = rock
		}
	}

	return cavern
}

type cavern [][]tile

type tile uint8

const (
	air tile = iota
	rock
	sand
)

type coord struct {
	x, y int
}

func parseCoord(s string) coord {
	prev := strings.SplitN(s, ",", 2)
	x, _ := strconv.Atoi(prev[0])
	y, _ := strconv.Atoi(prev[1])
	return coord{
		x: x,
		y: y,
	}
}

func (c coord) toward(o coord) (coord, bool) {
	switch {
	case c.x < o.x:
		return coord{x: c.x + 1, y: c.y}, true
	case c.x > o.x:
		return coord{x: c.x - 1, y: c.y}, true
	case c.y < o.y:
		return coord{x: c.x, y: c.y + 1}, true
	case c.y > o.y:
		return coord{x: c.x, y: c.y - 1}, true
	default:
		return c, false
	}
}

func (c coord) down() coord {
	return coord{x: c.x, y: c.y + 1}
}

func (c coord) diagonallyLeft() coord {
	return coord{x: c.x - 1, y: c.y + 1}
}

func (c coord) diagonallyRight() coord {
	return coord{x: c.x + 1, y: c.y + 1}
}

func (c *cavern) tile(coord coord) tile {
	return (*c)[coord.y][coord.x]
}

func (c *cavern) addSandUntilOneLandsInTheAbyss() int {
	count := 0
	cameToRest := true
	for cameToRest {
		cameToRest = c.addSand()
		if cameToRest {
			count++
		}
	}
	return count
}

func (c *cavern) addSandUntilSourceIsCovered() int {
	count := 0
	sourceIsNotBlocked := true
	for sourceIsNotBlocked {
		c.addSand()
		sourceIsNotBlocked = c.tile(coord{x: 500, y: 0}) == air
		count++
	}
	return count
}

func (c *cavern) addSand() bool {
	curPos := coord{x: 500, y: 0}
	falling := true
	for falling {
		if curPos.y >= 200-1 {
			// Into the abyss.
			return false
		}
		canFall := false
		for _, pos := range []coord{curPos.down(), curPos.diagonallyLeft(), curPos.diagonallyRight()} {
			if c.tile(pos) == air {
				curPos = pos
				canFall = true
				break
			}
		}
		if !canFall {
			break
		}
	}
	// Sand came to rest at curPos.
	(*c)[curPos.y][curPos.x] = sand
	return true
}

func (c *cavern) visualize(upperLeft, lowerRight coord) string {
	var b strings.Builder
	for _, row := range (*c)[upperLeft.y : lowerRight.y+1] {
		for x, tile := range row[upperLeft.x : lowerRight.x+1] {
			switch tile {
			case air:
				b.WriteRune('.')
			case rock:
				b.WriteRune('#')
			case sand:
				b.WriteRune('o')
			}
			if x+upperLeft.x == lowerRight.x {
				b.WriteRune('\n')
			}
		}
	}
	return b.String()
}
