package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input, err := os.ReadFile("go/day15/input.txt")
	if err != nil {
		panic(err)
	}
	sensorMap := parseSensorMap(string(input))
	fmt.Printf("Part 1. There are %d positions that cannot contain a beacon at y=2000000.\n", sensorMap.cannotContainBeaconCount(2_000_000))
	fmt.Printf("Part 2. The distress beacon has a tuning frequency of %d.\n", sensorMap.findDistressBeacon())
}

func parseSensorMap(input string) sensorMap {
	sensorMap := make(sensorMap)
	re := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, lineMatch := range matches {
		sensorX, _ := strconv.Atoi(lineMatch[1])
		sensorY, _ := strconv.Atoi(lineMatch[2])
		sensorCoord := coord{x: sensorX, y: sensorY}
		beaconX, _ := strconv.Atoi(lineMatch[3])
		beaconY, _ := strconv.Atoi(lineMatch[4])
		beaconCoord := coord{x: beaconX, y: beaconY}
		sensorMap[sensorCoord] = sensorInfo{
			loc:      sensorCoord,
			beacon:   beaconCoord,
			distance: sensorCoord.manhattanDistance(beaconCoord),
		}
	}
	return sensorMap
}

type sensorMap map[coord]sensorInfo

type coord struct {
	x, y int
}

func (c coord) manhattanDistance(o coord) uint {
	return abs(c.x-o.x) + abs(c.y-o.y)
}

func abs(i int) uint {
	if i < 0 {
		return uint(-i)
	}
	return uint(i)
}

func (c coord) within(s sensorInfo) (bool, uint) {
	dist := s.loc.manhattanDistance(c)
	// Non-optimal formula but it gets the job done.
	return dist <= s.distance, s.distance - abs(s.loc.x-c.x) - abs(s.loc.y-c.y)
}

func (c coord) equals(o coord) bool {
	return c.x == o.x && c.y == o.y
}

type sensorInfo struct {
	loc      coord
	beacon   coord
	distance uint
}

type device uint8

const (
	none device = iota
	sensor
	beacon
)

func (s sensorMap) findDistressBeacon() int {
	for y := 0; y <= 4_000_000; y++ {
		for x := 0; x <= 4_000_000; x++ {
			searchCoord := coord{x: x, y: y}
			canContain, dist := s.canContainBeacon(searchCoord, true)
			if canContain {
				return x*4_000_000 + y
			}
			if dist > 0 {
				x += int(dist) - 1
			}
		}
	}
	return -1
}

func (s sensorMap) cannotContainBeaconCount(y int) int {
	leftMost := 9999999
	rightMost := -9999999
	for _, sensor := range s {
		leftMostReach := sensor.loc.x - int(sensor.distance)
		if leftMostReach < leftMost {
			leftMost = leftMostReach
		}
		rightMostReach := sensor.loc.x + int(sensor.distance)
		if rightMostReach > rightMost {
			rightMost = rightMostReach
		}
	}
	blockedCount := 0
	for x := leftMost; x <= rightMost; x++ {
		checkCoord := coord{x: x, y: y}
		canContain, _ := s.canContainBeacon(checkCoord, false)
		if !canContain {
			blockedCount++
		}
	}
	return blockedCount
}

func (s sensorMap) canContainBeacon(c coord, ignoreScannerBeacons bool) (bool, uint) {
	for _, sensor := range s {
		if !ignoreScannerBeacons && c.equals(sensor.beacon) {
			return true, 0
		}
		isWithin, dist := c.within(sensor)
		if isWithin {
			return false, dist
		}
	}

	return true, 0
}
