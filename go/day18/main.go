package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("go/day18/input.txt")
	if err != nil {
		panic(err)
	}
	scan := parseScan(string(input))
	fmt.Printf("Part 1. The total surface area is %d\n", scan.surfaceArea(false))
	fmt.Printf("Part 2. The total exterior surface area is %d\n", scan.surfaceArea(true))
}

func parseScan(input string) scan {
	lava := make(map[coord]struct{})
	for _, line := range strings.Split(input, "\r\n") {
		coords := strings.SplitN(line, ",", 3)
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])
		coord := coord{x: x, y: y, z: z}
		lava[coord] = struct{}{}
	}
	return scan{
		lava:            lava,
		exteriorBubbles: make(map[coord]struct{}),
	}
}

type scan struct {
	lava            map[coord]struct{}
	exteriorBubbles map[coord]struct{}
}

type coord struct {
	x, y, z int
}

func (s scan) surfaceArea(onlyExterior bool) int {
	surfaceArea := 0
	for coord := range s.lava {
		surfaceArea += s.coordSurfaceArea(coord, onlyExterior)
	}
	return surfaceArea
}

func (s scan) coordSurfaceArea(c coord, onlyExterior bool) int {
	surfaceArea := 0
	for _, neighbor := range c.neighbors() {
		_, exists := s.lava[neighbor]
		if !exists && (!onlyExterior || s.isExteriorAirBubble(neighbor)) {
			surfaceArea++
		}
	}
	return surfaceArea
}

func (c coord) neighbors() []coord {
	return []coord{
		{x: c.x + 1, y: c.y, z: c.z},
		{x: c.x - 1, y: c.y, z: c.z},
		{x: c.x, y: c.y + 1, z: c.z},
		{x: c.x, y: c.y - 1, z: c.z},
		{x: c.x, y: c.y, z: c.z + 1},
		{x: c.x, y: c.y, z: c.z - 1},
	}
}

func (s scan) findMin(attr func(coord) int) int {
	min := 9999
	for l := range s.lava {
		a := attr(l)
		if a < min {
			min = a
		}
	}
	return min
}

func (s scan) findMax(attr func(coord) int) int {
	min := 0
	for l := range s.lava {
		a := attr(l)
		if a > min {
			min = a
		}
	}
	return min
}

func (s scan) minX() int {
	return s.findMin(func(c coord) int { return c.x })
}

func (s scan) maxX() int {
	return s.findMax(func(c coord) int { return c.x })
}

func (s scan) minY() int {
	return s.findMin(func(c coord) int { return c.y })
}

func (s scan) maxY() int {
	return s.findMax(func(c coord) int { return c.y })
}

func (s scan) minZ() int {
	return s.findMin(func(c coord) int { return c.z })
}

func (s scan) maxZ() int {
	return s.findMax(func(c coord) int { return c.z })
}

func (s scan) isExteriorAirBubble(c coord) bool {
	_, isExterior := s.exteriorBubbles[c]
	if isExterior {
		return true
	}
	return s.traverse(c, make(map[coord]struct{}), s.minX(), s.maxX(), s.minY(), s.maxY(), s.minZ(), s.maxZ())
}

func (s scan) traverse(c coord, visited map[coord]struct{}, minX, maxX, minY, maxY, minZ, maxZ int) bool {
	if c.x < minX || c.x > maxX || c.y < minY || c.y > maxY || c.z < minZ || c.z > maxZ {
		s.exteriorBubbles[c] = struct{}{}
		return true
	}
	visited[c] = struct{}{}

	for _, neighbor := range c.neighbors() {
		_, isLava := s.lava[neighbor]
		if isLava {
			continue
		}
		_, isExterior := s.exteriorBubbles[neighbor]
		_, hasVisited := visited[neighbor]
		if isExterior || !hasVisited && s.traverse(neighbor, visited, minX, maxX, minY, maxY, minZ, maxZ) {
			s.exteriorBubbles[c] = struct{}{}
			return true
		}
	}
	return false
}
