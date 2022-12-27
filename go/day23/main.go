package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day23/input.txt")
	if err != nil {
		panic(err)
	}
	grove := parseGrove(string(input))
	grove.stepN(10)
	fmt.Printf("Part 1. After 10 rounds there are %d empty spaces in the smallest possible rectangle\n", grove.emptySpacesInRectangle(grove.smallestRectangle()))
	grove = parseGrove(string(input))
	fmt.Printf("Part 2. Steps until no change: %d\n", grove.stepUntilNoChanges())
}

func parseGrove(input string) grove {
	grid := make(map[coord]struct{})
	for y, line := range strings.Split(input, "\r\n") {
		for x, r := range line {
			if r == '#' {
				coord := coord{x: x, y: y}
				grid[coord] = struct{}{}
			}
		}
	}
	return grove{
		grid:       grid,
		directions: [4]direction{north, south, west, east},
	}
}

type grove struct {
	grid       map[coord]struct{}
	directions [4]direction
}

type coord struct {
	x, y int
}

type direction int

const (
	invalid direction = iota
	north
	south
	west
	east
)

func (d direction) String() string {
	switch d {
	case north:
		return "north"
	case south:
		return "south"
	case west:
		return "west"
	case east:
		return "east"
	}
	return "invalid"
}

func (g *grove) stepUntilNoChanges() int {
	count := 0
	changed := true
	for changed {
		changed = g.step()
		count++
	}
	return count
}

func (g *grove) stepN(n int) {
	for i := 0; i < n; i++ {
		g.step()
	}
}

func (g *grove) step() bool {
	proposals := make(map[coord][]coord, len(g.grid))
	for c := range g.grid {
		var firstDirection direction
		canGoAllDirections := true
		for _, direction := range g.directions {
			d := direction
			canGo := g.canGo(c, d)
			if canGo && firstDirection == invalid {
				firstDirection = d
			}
			if !canGo {
				canGoAllDirections = false
			}
		}
		if canGoAllDirections || firstDirection == invalid {
			continue
		}
		var proposalCoord coord
		switch firstDirection {
		case north:
			proposalCoord = coord{x: c.x, y: c.y - 1}
		case south:
			proposalCoord = coord{x: c.x, y: c.y + 1}
		case west:
			proposalCoord = coord{x: c.x - 1, y: c.y}
		case east:
			proposalCoord = coord{x: c.x + 1, y: c.y}
		}
		proposals[proposalCoord] = append(proposals[proposalCoord], c)
	}

	didAMove := false
	for proposal, sources := range proposals {
		if len(sources) != 1 {
			continue
		}

		didAMove = true
		g.grid[proposal] = struct{}{}
		delete(g.grid, sources[0])
	}

	d := g.directions
	g.directions = [4]direction{d[1], d[2], d[3], d[0]}

	return didAMove
}

func (g *grove) has(c coord) bool {
	_, occupied := g.grid[c]
	return occupied
}

func (g *grove) canGo(c coord, d direction) bool {
	switch d {
	case north:
		return !g.has(coord{x: c.x, y: c.y - 1}) && !g.has(coord{x: c.x - 1, y: c.y - 1}) && !g.has(coord{x: c.x + 1, y: c.y - 1})
	case south:
		return !g.has(coord{x: c.x, y: c.y + 1}) && !g.has(coord{x: c.x - 1, y: c.y + 1}) && !g.has(coord{x: c.x + 1, y: c.y + 1})
	case west:
		return !g.has(coord{x: c.x - 1, y: c.y}) && !g.has(coord{x: c.x - 1, y: c.y + 1}) && !g.has(coord{x: c.x - 1, y: c.y - 1})
	case east:
		return !g.has(coord{x: c.x + 1, y: c.y}) && !g.has(coord{x: c.x + 1, y: c.y + 1}) && !g.has(coord{x: c.x + 1, y: c.y - 1})
	}
	panic("unknown direction")
}

func (g *grove) smallestRectangle() (minX, maxX, minY, maxY int) {
	maxX = -99999999999
	maxY = -99999999999
	minX = 99999999999
	minY = 99999999999
	for c := range g.grid {
		if c.x < minX {
			minX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
		if c.x > maxX {
			maxX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
	}
	return
}

func (g *grove) emptySpacesInRectangle(x1, x2, y1, y2 int) int {
	count := 0
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			_, occupied := g.grid[coord{x: x, y: y}]
			if !occupied {
				count++
			}
		}
	}
	return count
}

func (g *grove) String() string {
	var b strings.Builder
	x1, x2, y1, y2 := g.smallestRectangle()
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			_, occupied := g.grid[coord{x: x, y: y}]
			if occupied {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}
