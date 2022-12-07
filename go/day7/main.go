package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("go/day7/input.txt")
	root := parseDirectories(string(input))
	var sizes []uint64
	sizes, _ = root.allDirSizes(sizes)
	var smallest uint64 = 9999999999
	var totalSizeUnder100k uint64
	const (
		totalMemory    = 70000000
		requiredMemory = 30000000
	)
	usedMemory := root.size()
	freeMemory := totalMemory - usedMemory
	moreSpaceNeeded := requiredMemory - freeMemory
	for _, s := range sizes {
		if s < 100_000 {
			totalSizeUnder100k += s
		}
		if s > moreSpaceNeeded && s < smallest {
			smallest = s
		}
	}
	fmt.Printf("Part 1. The sum of all directory sizes under 100k is %d\n", totalSizeUnder100k)
	fmt.Printf("Part 2. The smallest directory you could delete to get the required space is %d\n", smallest)
}

func parseDirectories(input string) directory {
	root := directory{name: "/"}
	currentDirectory := &root
	previousDirectories := []directory{}
	for _, line := range strings.Split(input, "\r\n") {
		if strings.HasPrefix(line, "$ cd") {
			target := line[len("$ cd "):]
			switch target {
			case "/":
				currentDirectory = &root
				previousDirectories = []directory{}
			case "..":
				currentDirectory = &previousDirectories[len(previousDirectories)-1]
				// Pop last previous directory.
				if len(previousDirectories) > 1 {
					previousDirectories = previousDirectories[:len(previousDirectories)-1]
				} else {
					previousDirectories = []directory{}
				}
			default:
				var subDirectory *directory
				for _, sub := range currentDirectory.subdirectories {
					if sub.name == target {
						subDirectory = sub
						break
					}
				}
				previousDirectories = append(previousDirectories, *currentDirectory)
				currentDirectory = subDirectory
			}
		} else if line[0] != '$' {
			if strings.HasPrefix(line, "dir") {
				dir := &directory{
					name: line[len("dir "):],
				}
				currentDirectory.subdirectories = append(currentDirectory.subdirectories, dir)
			} else { // File.
				words := strings.SplitN(line, " ", 2)
				size, _ := strconv.ParseUint(words[0], 10, 32)
				file := file{
					name: words[1],
					size: size,
				}
				currentDirectory.files = append(currentDirectory.files, &file)
			}
		}
	}
	return root
}

type directory struct {
	name           string
	subdirectories []*directory
	files          []*file
}

type file struct {
	name string
	size uint64
}

func (d directory) size() uint64 {
	var sum uint64
	for _, f := range d.files {
		sum += f.size
	}
	for _, sub := range d.subdirectories {
		sum += sub.size()
	}

	return sum
}

func (d directory) allDirSizes(sizes []uint64) ([]uint64, uint64) {
	var sum uint64
	for _, f := range d.files {
		sum += f.size
	}
	for _, sub := range d.subdirectories {
		var subSum uint64
		sizes, subSum = sub.allDirSizes(sizes)
		sum += subSum
	}

	sizes = append(sizes, sum)
	return sizes, sum
}
