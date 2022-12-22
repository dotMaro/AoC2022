package main

import (
	"testing"
)

const example = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`

func Test_blueprint_largestNumberOfOpenedGeodesPossible(t *testing.T) {
	blueprints := parseBlueprints(example)
	res := blueprints[0].largestNumberOfOpenedGeodesPossible(24)
	if res != 9 {
		t.Errorf("Should return 9, not %d", res)
	}
	res = blueprints[1].largestNumberOfOpenedGeodesPossible(24)
	if res != 12 {
		t.Errorf("Should return 12, not %d", res)
	}
}
