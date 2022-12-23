package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	input, err := os.ReadFile("go/day20/input.txt")
	if err != nil {
		panic(err)
	}
	file := parseFile(string(input))
	file.decryptN(1)
	fmt.Printf("Part 1. The sum of the coordinates is %d\n", file.sumOfCoordinates())
	file = parseFile(string(input))
	timeStart := time.Now()
	file.applyDecryptionKey(811589153)
	file.decryptN(10)
	fmt.Printf("Part 2 (%s). The sum of the coordinates after applying the decryption key is %d\n", time.Since(timeStart), file.sumOfCoordinates())
}

func parseFile(input string) file {
	lines := strings.Split(input, "\r\n")
	file := make(file, 0, len(lines))
	for _, line := range lines {
		val, _ := strconv.Atoi(line)
		file = append(file, val)
	}
	return file
}

type file []int

func (f file) applyDecryptionKey(key int) {
	for i := range f {
		f[i] *= key
	}
}

func (f file) decryptN(n int) {
	indexOffset := make([]int, len(f))
	for i := 0; i < n; i++ {
		indexOffset = f.decrypt(indexOffset)
	}
}

func (f file) decrypt(indexOffset []int) []int {
	for i, offset := range indexOffset {
		indexBeforeMove := i + offset
		valToMove := f[indexBeforeMove]

		// Pop.
		f = append(f[:indexBeforeMove], f[indexBeforeMove+1:]...)
		newIndex := (indexBeforeMove + valToMove) % len(f)
		for newIndex < 0 {
			newIndex += len(f)
		}
		tail := make([]int, len(f[newIndex:]))
		copy(tail, f[newIndex:])
		f = append(f[:newIndex], valToMove)
		f = append(f, tail...)

		// Adjust indexOffset.
		for s, offset2 := range indexOffset {
			if s+offset2 > indexBeforeMove && s+offset2 <= newIndex {
				indexOffset[s]--
				if s+indexOffset[s] < 0 {
					indexOffset[s] += len(f)
				}
			} else if newIndex < indexBeforeMove && s+offset2 < indexBeforeMove && s+offset2 >= newIndex {
				indexOffset[s] = (indexOffset[s] + 1) % len(f)
			}
		}
		indexOffset[i] = newIndex - i
	}
	return indexOffset
}

func (f file) sumOfCoordinates() int {
	zeroIndex := 0
	for i, v := range f {
		if v == 0 {
			zeroIndex = i
			break
		}
	}

	return f[(zeroIndex+1000)%len(f)] + f[(zeroIndex+2000)%len(f)] + f[(zeroIndex+3000)%len(f)]
}
