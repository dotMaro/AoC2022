package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, _ := os.ReadFile("go/day12/input.txt")
	heightMap := parse(string(input))
	fmt.Printf("Part 1. The shortest path is %d steps.\n", heightMap.traverse(heightMap.start, func(c coord) bool { return heightMap.height(c) == 'E' }, heightMap.canGo))
	fmt.Printf("Part 2. The shortest path is %d steps.\n", heightMap.traverse(heightMap.end, func(c coord) bool { return heightMap.height(c) == 'a' }, heightMap.canGoReverse))
}

func parse(input string) heightMap {
	var start, end coord
	var graph [][]byte
	for y, line := range strings.Split(input, "\r\n") {
		for x, c := range line {
			switch c {
			case 'S':
				start = coord{x: x, y: y}
			case 'E':
				end = coord{x: x, y: y}
			}
		}
		graph = append(graph, []byte(line))
	}
	return heightMap{
		graph: graph,
		start: start,
		end:   end,
	}
}

type heightMap struct {
	graph      [][]byte
	start, end coord
}

type coord struct {
	x, y int
}

func (m *heightMap) neighbors(c coord) []coord {
	var neighbors []coord
	if c.x > 0 {
		neighbors = append(neighbors, coord{x: c.x - 1, y: c.y})
	}
	if c.x < len(m.graph[0])-1 {
		neighbors = append(neighbors, coord{x: c.x + 1, y: c.y})
	}
	if c.y > 0 {
		neighbors = append(neighbors, coord{x: c.x, y: c.y - 1})
	}
	if c.y < len(m.graph)-1 {
		neighbors = append(neighbors, coord{x: c.x, y: c.y + 1})
	}
	return neighbors
}

func (m *heightMap) canGo(source, target coord) bool {
	sourceHeight := m.height(source)
	targetHeight := m.height(target)
	if sourceHeight == 'S' {
		sourceHeight = 'a'
	}
	if targetHeight == 'E' {
		targetHeight = 'z'
	}
	return int(targetHeight) == int(sourceHeight)+1 || int(targetHeight) <= int(sourceHeight)
}

func (m *heightMap) canGoReverse(source, target coord) bool {
	return m.canGo(target, source)
}

func (m *heightMap) height(c coord) byte {
	return m.graph[c.y][c.x]
}

func (m *heightMap) traverse(from coord, toCondition func(coord) bool, canGo func(coord, coord) bool) uint {
	distance := make(map[coord]uint)
	previous := make(map[coord]coord)
	visited := make(map[coord]struct{})
	q := make(map[coord]struct{})

	for y := range m.graph {
		for x := range m.graph[y] {
			coord := coord{x: x, y: y}
			distance[coord] = 999999999999
			q[coord] = struct{}{}
		}
	}
	distance[from] = 0

	for len(q) > 0 {
		var node coord
		var min uint = 999999999999
		for candidate := range q {
			dist, hasDist := distance[candidate]
			if hasDist && dist < min {
				node = candidate
				min = dist
			}
		}
		if toCondition(node) {
			return distance[node]
		}
		delete(q, node)
		for _, neighbor := range m.neighbors(node) {
			_, hasVisited := visited[neighbor]
			if hasVisited || !canGo(node, neighbor) {
				continue
			}
			visited[neighbor] = struct{}{}

			alt := distance[node] + 1
			if alt < distance[neighbor] {
				distance[neighbor] = alt
				previous[neighbor] = node
			}
		}
	}

	return 0
}
