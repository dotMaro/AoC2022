package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// A whole lot of spaghetti for this one. In the end I remade it to use memoization and pre-filling the opened valve set with all possible configurations to calculate the highest possible flow.
// Credit to this guy https://www.youtube.com/watch?v=bLMj50cpOug

func main() {
	input, _ := os.ReadFile("go/day16/input.txt")
	valveMap := parseValves(string(input))
	fmt.Printf("Part 1. The highest amount of pressure that can be released is %d", valveMap.findOptimalOrderForMaximumPressureRelease())
	fmt.Printf("Part 2. The highest amount of pressure that can be released is %d", valveMap.findOptimalOrderForMaximumPressureReleaseWithElephant())
}

func parseValves(input string) valveMap {
	valveMap := make(valveMap)
	re := regexp.MustCompile(`Valve (.+) has flow rate=(\d+); tunnel(s?) lead(s?) to valve(s?) (.+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, lineMatch := range matches {
		valveName := lineMatch[1]
		flowRate, _ := strconv.Atoi(lineMatch[2])
		tunnelStrings := strings.Split(lineMatch[6], ", ")
		tunnels := make([]string, 0, len(tunnelStrings))
		for _, t := range tunnelStrings {
			tunnels = append(tunnels, strings.TrimSuffix(t, "\r"))
		}
		valveMap[valveName] = valve{
			name:     valveName,
			flowRate: flowRate,
			leadsTo:  tunnels,
		}
	}
	return valveMap
}

type valveMap map[string]valve

type valve struct {
	name     string
	flowRate int
	leadsTo  []string
}

func (v valveMap) findOptimalOrderForMaximumPressureRelease() int {
	valveIndices := make(map[string]int)
	i := 0
	for _, valve := range v {
		if valve.flowRate == 0 {
			continue
		}
		valveIndices[valve.name] = i
		i++
	}
	return v.dfs(30, v["AA"], 0, v.shortestPathsBetweenValvesWithNonZeroFlowRate(), valveIndices, make(cache))
}

func (v valveMap) findOptimalOrderForMaximumPressureReleaseWithElephant() int {
	valveIndices := make(map[string]int)
	i := 0
	for _, valve := range v {
		if valve.flowRate == 0 {
			continue
		}
		valveIndices[valve.name] = i
		i++
	}
	allOpened := (1 << len(valveIndices)) - 1
	shortestPaths := v.shortestPathsBetweenValvesWithNonZeroFlowRate()
	max := 0
	cache := make(cache)
	// /2 to remove duplicated checks.
	for i := 0; i <= allOpened/2; i++ {
		res := v.dfs(26, v["AA"], i, shortestPaths, valveIndices, cache) + v.dfs(26, v["AA"], allOpened^i, shortestPaths, valveIndices, cache)
		if res > max {
			max = res
		}
	}
	return max
}

func (v valveMap) release(minute int, pressureReleased int, pressureReleasedPerMin int, currentValve valve, openValves map[string]struct{}, shortestPaths map[string][][]string, withElephant bool) (int, []string) {
	if minute == 1 {
		released := 0
		var trace []string
		if withElephant {
			released, trace = v.release(26, 0, 0, v["AA"], openValves, shortestPaths, false)
		}
		return pressureReleased + pressureReleasedPerMin + released, append(trace, []string{"<-ELEPHANTO", currentValve.name}...)
	}

	highest := 0
	var highestTrace []string
	// Only try going to valves with a flow rate. Find cheapest path to each of them.
	pathFound := false
	for _, path := range shortestPaths[currentValve.name] {
		_, alreadyOpen := openValves[path[0]]
		if alreadyOpen || minute-len(path) <= 0 {
			continue
		}
		pathFound = true
		target := v[path[0]]
		copiedOpenValves := make(map[string]struct{}, len(openValves))
		for key := range openValves {
			copiedOpenValves[key] = struct{}{}
		}
		// released, trace := v.release(minute-len(path)-1, pressureReleased+pressureReleasedPerMin*(len(path)-1), pressureReleasedPerMin, target, copiedOpenValves, shortestPaths)

		// len(path) includes the current valve, but we want to +1 for the time it takes to open the valve anyway.
		var released int
		var trace []string
		copiedOpenValves[path[0]] = struct{}{}
		released, trace = v.release(minute-len(path), pressureReleased+pressureReleasedPerMin*len(path), pressureReleasedPerMin+target.flowRate, target, copiedOpenValves, shortestPaths, withElephant)
		if released > highest {
			highest = released
			highestTrace = trace
		}
	}
	if !pathFound {
		// No more paths fit but add pressureReleasedPerMin.
		highest = pressureReleased + minute*pressureReleasedPerMin
	}

	return highest, append(highestTrace, currentValve.name)
}

// This is such a mess, but it actually got pretty close to the correct answer (2580 vs 2602).
func (v valveMap) releaseWithElephant(minute int, pressureReleased int, pressureReleasedPerMin int, santaValve valve, elephantValve valve, santaPath []string, elephantPath []string, openValves map[string]struct{}, shortestPaths map[string][][]string) (int, []string) {
	if minute == 1 {
		return pressureReleased + pressureReleasedPerMin, []string{santaValve.name}
	}

	highest := 0
	var highestTrace []string
	paths := v.targets(minute, santaValve, elephantValve, santaPath, elephantPath, openValves, shortestPaths)
	for _, path := range paths {
		var released int
		var trace []string
		if len(path.santa) < len(path.elephant) && len(path.santa) > 0 {
			target := v[path.santa[0]]
			released, trace = v.releaseWithElephant(minute-len(path.santa), pressureReleased+pressureReleasedPerMin*len(path.santa), pressureReleasedPerMin+target.flowRate, target, elephantValve, nil, path.elephant[:len(path.elephant)-len(path.santa)], openValves, shortestPaths)
		} else if len(path.elephant) < len(path.santa) && len(path.elephant) > 0 {
			target := v[path.elephant[0]]
			released, trace = v.releaseWithElephant(minute-len(path.elephant), pressureReleased+pressureReleasedPerMin*len(path.elephant), pressureReleasedPerMin+target.flowRate, santaValve, target, path.santa[:len(path.santa)-len(path.elephant)], nil, openValves, shortestPaths)
		} else {
			elephantTarget := v[path.elephant[0]]
			santaTarget := v[path.santa[0]]
			released, trace = v.releaseWithElephant(minute-len(path.elephant), pressureReleased+pressureReleasedPerMin*len(path.elephant), pressureReleasedPerMin+santaTarget.flowRate+elephantTarget.flowRate, santaTarget, elephantTarget, nil, nil, openValves, shortestPaths)
		}
		if released > highest {
			highest = released
			highestTrace = trace
		}
	}
	if len(paths) == 0 {
		// No more paths fit but add pressureReleasedPerMin.
		highest = pressureReleased + minute*pressureReleasedPerMin
	}

	return highest, append(highestTrace, santaValve.name)
}

type target struct {
	santa    []string
	elephant []string
}

func (v valveMap) targets(minute int, santaValve valve, elephantValve valve, currentSantaPath []string, currentElephantPath []string, openValves map[string]struct{}, shortestPaths map[string][][]string) []target {
	var targets []target

	// Handle case when there's only one path left to do. <<<<<<<<<<<<<<<<<<<<<<
	if len(currentSantaPath) == 0 && len(currentElephantPath) == 0 {
		for _, santaPath := range shortestPaths[santaValve.name] {
			_, alreadyOpen := openValves[santaPath[0]]
			if alreadyOpen || minute-len(santaPath) <= 0 {
				continue
			}
			// foundElephantPath := false
			for _, elephantPath := range shortestPaths[elephantValve.name] {
				// if elephantPath[0] == santaPath[0] {
				// 	continue
				// }
				_, alreadyOpen := openValves[elephantPath[0]]
				if alreadyOpen || minute-len(elephantPath) <= 0 {
					continue
				}
				// foundElephantPath = true
				copiedOpenValves := make(map[string]struct{}, len(openValves))
				for key := range openValves {
					copiedOpenValves[key] = struct{}{}
				}
				target := target{
					santa:    santaPath,
					elephant: elephantPath,
				}
				targets = append(targets, target)
			}
			// if !foundElephantPath {
			// 	target := target{
			// 		santa:    santaPath,
			// 		elephant: nil,
			// 	}
			// 	targets = append(targets, target)
			// }
		}
	} else if len(currentSantaPath) > 0 {
		for _, elephantPath := range shortestPaths[elephantValve.name] {
			_, alreadyOpen := openValves[elephantPath[0]]
			if alreadyOpen || minute-len(elephantPath) <= 0 {
				continue
			}
			target := target{
				santa:    currentSantaPath,
				elephant: elephantPath,
			}
			targets = append(targets, target)
		}
	} else if len(currentElephantPath) > 0 {
		for _, santaPath := range shortestPaths[santaValve.name] {
			_, alreadyOpen := openValves[santaPath[0]]
			if alreadyOpen || minute-len(santaPath) <= 0 {
				continue
			}
			target := target{
				santa:    santaPath,
				elephant: currentElephantPath,
			}
			targets = append(targets, target)
		}
	} else {
		panic("wtf")
	}

	return targets
}

// from -> paths to different targets (with flow rate >0)
func (v valveMap) shortestPathsBetweenValvesWithNonZeroFlowRate() map[string][][]string {
	shortestPaths := make(map[string][][]string)
	for _, valve1 := range v {
		if valve1.flowRate == 0 && valve1.name != "AA" {
			continue
		}
		for _, valve2 := range v {
			if valve2.name == valve1.name || valve2.flowRate == 0 {
				continue
			}
			shortestPaths[valve1.name] = append(shortestPaths[valve1.name], v.traverse(valve1, valve2, make(map[string]struct{})))
		}
	}
	return shortestPaths
}

func (v valveMap) traverse(current, target valve, visited map[string]struct{}) []string {
	visited[current.name] = struct{}{}
	shortestPathLen := 999
	var shortestPath []string
	for _, valve := range current.leadsTo {
		if valve == target.name {
			return []string{target.name, current.name}
		}
		_, alreadyVisited := visited[valve]
		if !alreadyVisited {
			copiedVisited := make(map[string]struct{})
			for k, val := range visited {
				copiedVisited[k] = val
			}
			path := v.traverse(v[valve], target, copiedVisited)
			if len(path) < shortestPathLen && path[0] == target.name {
				shortestPathLen = len(path)
				shortestPath = path
			}
		}
	}
	return append(shortestPath, current.name)
}

type cache map[params]int

type params struct {
	time      int
	valveName string
	opened    int
}

func (v valveMap) dfs(time int, valve valve, opened int, edges map[string][][]string, valveIndexes map[string]int, cache cache) int {
	params := params{
		time:      time,
		valveName: valve.name,
		opened:    opened,
	}
	cacheRes, inCache := cache[params]
	if inCache {
		return cacheRes
	}

	maxVal := 0

	for _, neighbors := range edges[valve.name] {
		// -1 since the neighbor list includes the current valve.
		// +1 since we're opening the target valve.
		remainingTime := time - (len(neighbors) - 1 + 1)
		if remainingTime <= 0 {
			continue
		}
		bit := 1 << valveIndexes[neighbors[0]]
		if opened&bit != 0 {
			continue
		}
		res := v.dfs(remainingTime, v[neighbors[0]], opened|bit, edges, valveIndexes, cache) + v[neighbors[0]].flowRate*remainingTime
		if res > maxVal {
			maxVal = res
		}
	}

	cache[params] = maxVal
	return maxVal
}
