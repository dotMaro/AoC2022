package main

import (
	"strings"
	"testing"
)

const example = `1
2
-3
3
-2
0
4`

const custom = `-1
4
0
10
-7
2`

func Test_example(t *testing.T) {
	file := parseFile(strings.ReplaceAll(example, "\n", "\r\n"))
	file.applyDecryptionKey(811589153)
	file.decryptN(10)
	sum := file.sumOfCoordinates()
	if sum != 1623178306 {
		t.Errorf("Should return 1623178306, not %d", sum)
	}
	// t.Fail()
}
