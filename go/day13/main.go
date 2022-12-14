package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	dividerPacket1 = "[[2]]"
	dividerPacket2 = "[[6]]"
)

func main() {
	input, _ := os.ReadFile("go/day13/input.txt")

	pairs := parsePairs(string(input))
	fmt.Printf("Part 1. Sum of pair indexes with correct packet orders: %d\n", sumOfPairsWithCorrectPacketOrder(pairs))

	packets := parsePackets(string(input))
	packets = append(packets, parsePacket(dividerPacket1), parsePacket(dividerPacket2))
	fmt.Printf("Part 2. Product of divider packets is %d\n", productOfDividers(packets))
}

func productOfDividers(packets []packet) int {
	sort.Slice(packets, func(i, j int) bool {
		isRight, _ := packets[i].isRightOrder(packets[j])
		return isRight
	})
	dividerSum := 1
	for i, packet := range packets {
		if packet.s == dividerPacket1 || packet.s == dividerPacket2 {
			dividerSum *= i + 1
		}
	}
	return dividerSum
}

func sumOfPairsWithCorrectPacketOrder(pairs []pair) int {
	sum := 0
	for i, pair := range pairs {
		isRight, _ := pair.left.isRightOrder(pair.right)
		if isRight {
			sum += i + 1
		}
	}
	return sum
}

func parsePairs(input string) []pair {
	var pairs []pair
	for _, strPair := range strings.Split(string(input), "\r\n\r\n") {
		lines := strings.SplitN(strPair, "\r\n", 2)
		left := parsePacket(lines[0])
		right := parsePacket(lines[1])
		pair := pair{
			left:  left,
			right: right,
		}
		pairs = append(pairs, pair)
	}
	return pairs
}

type pair struct {
	left  packet
	right packet
}

type packet struct {
	data []any
	s    string
}

func parsePackets(input string) []packet {
	var packets []packet
	for _, line := range strings.Split(input, "\r\n") {
		if line == "" {
			continue
		}
		packets = append(packets, parsePacket(line))
	}
	return packets
}

func parsePacket(s string) packet {
	var data []any
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		panic(err)
	}
	return packet{
		data: data,
		s:    s,
	}
}

// isRight and continue
func (p packet) isRightOrder(other packet) (bool, bool) {
	for i, leftListValue := range p.data {
		if i > len(other.data)-1 {
			return false, false
		}
		rightListValue := other.data[i]
		leftListValueNumber, leftValueIsNumber := leftListValue.(float64)
		rightListValueNumber, rightValueIsNumber := rightListValue.(float64)
		if leftValueIsNumber && rightValueIsNumber {
			if leftListValueNumber < rightListValueNumber {
				return true, false
			} else if leftListValueNumber > rightListValueNumber {
				return false, false
			} else if leftListValueNumber == rightListValueNumber {
				continue
			}
		} else if !leftValueIsNumber && !rightValueIsNumber {
			leftPacket := packet{data: leftListValue.([]any)}
			rightPacket := packet{data: rightListValue.([]any)}
			isRight, cont := leftPacket.isRightOrder(rightPacket)
			if cont {
				continue
			}
			return isRight, false
		} else { // Mixed.
			if leftValueIsNumber {
				leftPacket := packet{data: []any{leftListValueNumber}}
				rightPacket := packet{data: rightListValue.([]any)}
				isRight, cont := leftPacket.isRightOrder(rightPacket)
				if cont {
					continue
				}
				return isRight, false
			}
			leftPacket := packet{data: leftListValue.([]any)}
			rightPacket := packet{data: []any{rightListValueNumber}}
			isRight, cont := leftPacket.isRightOrder(rightPacket)
			if cont {
				continue
			}
			return isRight, false
		}
	}
	if len(p.data) > len(other.data) {
		return false, false
	} else if len(p.data) < len(other.data) {
		return true, false
	} else if len(p.data) == len(other.data) {
		return false, true
	}
	panic(fmt.Sprintf("Couldn't handle %s and %s", p.data, other.data))
}
