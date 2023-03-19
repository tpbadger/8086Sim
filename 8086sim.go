// 8086sim is a sim for the Intel 8086 cpu. Currently being developed as part of
// performance aware programming course by Casey M
// usage -> go run 8086sim.go <binary instruction file>

package main

import (
	"os"
)

func main() {
	s := os.Args[1]

	memory := Memory{}

	err := LoadFileIntoMemory(s, &memory)

	if err != nil {
		// cba to deal with thim imo
		panic(err)
	}

	DecodeInstructions(&memory)
}
