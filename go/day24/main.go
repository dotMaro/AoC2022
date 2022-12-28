package main

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	blizzardCache = make(map[state]struct{})
}

func main() {
	input, err := os.ReadFile("go/day24/input.txt")
	if err != nil {
		panic(err)
	}
	valley := parseValley(string(input))
	fmt.Printf("Part 1. The shortest time to reach the end is %d minutes\n", valley.findShortestPath(false))
	fmt.Printf("Part 2. The shortest time to reach the end, go back, then to the end again is %d minutes\n", valley.findShortestPath(true))
}

func parseValley(input string) valley {
	lines := strings.Split(input, "\r\n")
	valley := make(valley, len(lines))
	for y, line := range lines {
		valley[y] = make([][]direction, len(line))
		for x, r := range line {
			if r == '.' {
				continue
			}
			var d direction
			switch r {
			case '^':
				d = up
			case 'v':
				d = down
			case '>':
				d = right
			case '<':
				d = left
			case '#':
				d = permanent
			}
			valley[y][x] = append(valley[y][x], d)
		}
	}

	valley.fillBlizzardMap()
	return valley
}

type valley [][][]direction

type direction int

const (
	up direction = iota
	down
	right
	left
	permanent
)

type state struct {
	x, y int
	t    int
}

func (v valley) findShortestPath(doubleJourney bool) int {
	queue := []state{
		{
			x: 1,
			y: 0,
			t: 0,
		},
	}
	seen := make(map[state]struct{}, 1000)
	reachedGoal := false
	reachedStart := false

	for len(queue) > 0 {
		s := queue[0]
		if len(queue) > 0 {
			queue = queue[1:]
		} else {
			queue = []state{}
		}

		_, hasSeen := seen[s]
		if hasSeen {
			continue
		}
		seen[s] = struct{}{}

		if s.y == len(v)-1 {
			if !doubleJourney {
				return s.t
			}
			if !reachedGoal {
				reachedGoal = true
				seen = make(map[state]struct{}, 1000)
				queue = []state{}
			} else if reachedStart {
				return s.t
			}
		} else if reachedGoal && !reachedStart && s.y == 0 {
			reachedStart = true
			seen = make(map[state]struct{}, 1000)
			queue = []state{}
		}

		for _, ss := range []state{
			{x: s.x + 1, y: s.y, t: s.t + 1},
			{x: s.x, y: s.y + 1, t: s.t + 1},
			{x: s.x - 1, y: s.y, t: s.t + 1},
			{x: s.x, y: s.y - 1, t: s.t + 1},
			{x: s.x, y: s.y, t: s.t + 1},
		} {
			if ss.y < 0 || ss.y >= len(v) {
				continue
			}
			if !v.hasBlizzard(ss) {
				queue = append(queue, ss)
			}
		}
	}
	return -1
}

var blizzardCache map[state]struct{}

func (v valley) hasBlizzard(s state) bool {
	s.t %= (len(v) - 2) * (len(v[0]) - 2)
	_, hasBlizzard := blizzardCache[s]
	return hasBlizzard
}

func (v valley) fillBlizzardMap() {
	for t := 0; t <= (len(v)-2)*(len(v[0])-2)+1; t++ {
		for y := 0; y < len(v); y++ {
			for x := 0; x < len(v[0]); x++ {
				for _, direction := range v[y][x] {
					switch direction {
					case permanent:
						s := state{x: x, y: y, t: t}
						blizzardCache[s] = struct{}{}
					case right:
						s := state{x: 1 + (x-1+t)%(len(v[0])-2), y: y, t: t}
						blizzardCache[s] = struct{}{}
					case left:
						s := state{x: (x - 1 - t) % (len(v[0]) - 2), y: y, t: t}
						if s.x < 0 {
							s.x += len(v[0]) - 2
						}
						s.x++
						blizzardCache[s] = struct{}{}
					case up:
						s := state{x: x, y: (y - 1 - t) % (len(v) - 2), t: t}
						if s.y < 0 {
							s.y += len(v) - 2
						}
						s.y++
						blizzardCache[s] = struct{}{}
					case down:
						s := state{x: x, y: 1 + (y-1+t)%(len(v)-2), t: t}
						blizzardCache[s] = struct{}{}
					}
				}
			}
		}
		// s := state{x: 1, y: -1, t: t}
		// blizzardCache[s] = struct{}{}
	}
}

func (v valley) visualize(t int) string {
	var b strings.Builder
	for y := 0; y < len(v); y++ {
		for x := 0; x < len(v[0]); x++ {
			_, inCache := blizzardCache[state{x: x, y: y, t: t}]
			if inCache {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}
