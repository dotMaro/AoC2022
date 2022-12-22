package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input, err := os.ReadFile("go/day19/input.txt")
	if err != nil {
		panic(err)
	}
	blueprints := parseBlueprints(string(input))
	fmt.Printf("Part 1. The combined quality level is %d\n", totalQualityLevel(blueprints))
	fmt.Printf("Part 2. The product of the first three blueprints is %d\n", productOfHighestGeodes(blueprints[:3]))
}

func totalQualityLevel(blueprints []blueprint) int {
	sum := 0
	for i, b := range blueprints {
		sum += (i + 1) * int(b.largestNumberOfOpenedGeodesPossible(24))
	}
	return sum
}

func productOfHighestGeodes(blueprints []blueprint) int {
	product := 1
	for i, b := range blueprints {
		fmt.Printf("%d/%d\n", i+1, len(blueprints))
		product *= int(b.largestNumberOfOpenedGeodesPossible(32))
	}
	return product
}

func parseBlueprints(input string) []blueprint {
	re := regexp.MustCompile(`Blueprint \d+: Each ore robot costs (\d+) ore\. Each clay robot costs (\d+) ore\. Each obsidian robot costs (\d+) ore and (\d+) clay\. Each geode robot costs (\d+) ore and (\d+) obsidian\.`)
	matches := re.FindAllStringSubmatch(input, -1)
	blueprints := make([]blueprint, 0, len(matches))
	for _, subMatches := range matches {
		blueprint := blueprint{
			oreRobotOreCost:  unwrap(strconv.Atoi(subMatches[1])),
			clayRobotOreCost: unwrap(strconv.Atoi(subMatches[2])),
			obsidianRobotCost: struct {
				ore  uint8
				clay uint8
			}{
				ore:  unwrap(strconv.Atoi(subMatches[3])),
				clay: unwrap(strconv.Atoi(subMatches[4])),
			},
			geodeRobotCost: struct {
				ore      uint8
				obsidian uint8
			}{
				ore:      unwrap(strconv.Atoi(subMatches[5])),
				obsidian: unwrap(strconv.Atoi(subMatches[6])),
			},
		}
		blueprints = append(blueprints, blueprint)
	}
	return blueprints
}

func unwrap(v int, err error) uint8 {
	if err != nil {
		panic(err)
	}
	return uint8(v)
}

type blueprint struct {
	oreRobotOreCost   uint8
	clayRobotOreCost  uint8
	obsidianRobotCost struct {
		ore  uint8
		clay uint8
	}
	geodeRobotCost struct {
		ore      uint8
		obsidian uint8
	}
}

type signature struct {
	minute                                             uint8
	oreCount, clayCount, obsidianCount                 uint8
	oreRobots, clayRobots, obsidianRobots, geodeRobots uint8
}

func (b blueprint) largestNumberOfOpenedGeodesPossible(maxMinutes uint8) uint8 {
	return b.traverse(maxMinutes, 1, 0, 0, 0, 1, 0, 0, 0, make(map[signature]uint8))
}

func (b blueprint) traverse(maxMinutes uint8, minute uint8, oreCount, clayCount, obsidianCount uint8, oreRobots, clayRobots, obsidianRobots, geodeRobots uint8, highest map[signature]uint8) uint8 {
	if minute == maxMinutes {
		return geodeRobots
	}
	totalRobotCount := oreRobots + clayRobots + obsidianRobots + geodeRobots
	if maxMinutes > 30 && minute > 21 && (totalRobotCount < 8 || geodeRobots == 0) {
		return 0
	}

	sign := signature{
		minute:         minute,
		oreCount:       oreCount,
		clayCount:      clayCount,
		obsidianCount:  obsidianCount,
		oreRobots:      oreRobots,
		clayRobots:     clayRobots,
		obsidianRobots: obsidianRobots,
		geodeRobots:    geodeRobots,
	}
	h, seen := highest[sign]
	if seen {
		return h
	}

	var (
		oreRes      uint8
		clayRes     uint8
		obsidianRes uint8
		geodeRes    uint8
		nothingRes  uint8
	)
	if oreCount >= b.oreRobotOreCost {
		oreRes = b.traverse(maxMinutes, minute+1, oreCount+oreRobots-b.oreRobotOreCost, clayCount+clayRobots, obsidianCount+obsidianRobots, oreRobots+1, clayRobots, obsidianRobots, geodeRobots, highest)
	}
	if oreCount >= b.clayRobotOreCost {
		clayRes = b.traverse(maxMinutes, minute+1, oreCount+oreRobots-b.clayRobotOreCost, clayCount+clayRobots, obsidianCount+obsidianRobots, oreRobots, clayRobots+1, obsidianRobots, geodeRobots, highest)
	}
	if oreCount >= b.obsidianRobotCost.ore && clayCount >= b.obsidianRobotCost.clay {
		obsidianRes = b.traverse(maxMinutes, minute+1, oreCount+oreRobots-b.obsidianRobotCost.ore, clayCount+clayRobots-b.obsidianRobotCost.clay, obsidianCount+obsidianRobots, oreRobots, clayRobots, obsidianRobots+1, geodeRobots, highest)
	}
	if oreCount >= b.geodeRobotCost.ore && obsidianCount >= b.geodeRobotCost.obsidian {
		geodeRes = b.traverse(maxMinutes, minute+1, oreCount+oreRobots-b.geodeRobotCost.ore, clayCount+clayRobots, obsidianCount+obsidianRobots-b.geodeRobotCost.obsidian, oreRobots, clayRobots, obsidianRobots, geodeRobots+1, highest)
	}
	// Do not construct any robot.
	nothingRes = b.traverse(maxMinutes, minute+1, oreCount+oreRobots, clayCount+clayRobots, obsidianCount+obsidianRobots, oreRobots, clayRobots, obsidianRobots, geodeRobots, highest)

	highestRes := max(oreRes, clayRes, obsidianRes, geodeRes, nothingRes)
	highest[sign] = highestRes
	return highestRes + geodeRobots
}

func max(n ...uint8) uint8 {
	highest := uint8(0)
	for _, val := range n {
		if val > highest {
			highest = val
		}
	}
	return highest
}
