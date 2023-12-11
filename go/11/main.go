package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

type Pos [2]int
type Galaxies []Pos

func getGalaxies(in []byte) Galaxies {
	width := bytes.IndexByte(in, '\n') + 1
	galaxies := make([]Pos, 0, 1024)
	for i, c := range in {
		if c == '#' {
			galaxies = append(galaxies, Pos{i % width, i / width})
		}
	}
	return galaxies
}

func expand(galaxies Galaxies, factor int) {
	for dim := 0; dim < 2; dim++ {
		slices.SortFunc(galaxies, func(a, b Pos) int {
			return a[dim] - b[dim]
		})
		spaceOffset := 0
		for i, end := 0, len(galaxies)-1; i < end; i++ {
			dx := galaxies[i+1][dim] + spaceOffset - galaxies[i][dim]
			spaceOffset += max(dx-1, 0) * (factor - 1)
			galaxies[i+1][dim] += spaceOffset
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func solve(in []byte, scaleFactor int) int {
	galaxies := getGalaxies(in)
	expand(galaxies, scaleFactor)
	res := 0
	for i, end := 0, len(galaxies); i < end; i++ {
		for j := i + 1; j < end; j++ {
			for dim := 0; dim < 2; dim++ {
				res += abs(galaxies[i][dim] - galaxies[j][dim])
			}
		}
	}
	return res
}

func solve1(in []byte) int {
	return solve(in, 2)
}

func solve2(in []byte) int {
	return solve(in, 1000000)
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
