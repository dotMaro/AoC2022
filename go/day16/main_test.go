package main

import (
	"strings"
	"testing"
)

const example = `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`

func Test_valveMap_findOptimalOrderForMaximumPressureRelease(t *testing.T) {
	valveMap := parseValves(strings.ReplaceAll(example, "\n", "\r\n"))
	res := valveMap.findOptimalOrderForMaximumPressureRelease()
	if res != 1651 {
		t.Errorf("Should return 1651, but got %d", res)
	}
}

func Test_valveMap_findOptimalOrderForMaximumPressureReleaseWithElephant(t *testing.T) {
	valveMap := parseValves(strings.ReplaceAll(example, "\n", "\r\n"))
	res := valveMap.findOptimalOrderForMaximumPressureReleaseWithElephant()
	if res != 1707 {
		t.Errorf("Should return 1707, but got %d", res)
	}
}
