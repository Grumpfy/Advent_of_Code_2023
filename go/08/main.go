package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

type PuzzleData struct {
	dirs        []byte
	vertexIndex map[string]int
	childs      [][2]int
}

func parseInput(in []byte) PuzzleData {
	res := PuzzleData{}
	scanner := bufio.NewScanner(bytes.NewReader(in))
	// directions
	scanner.Scan()
	dirs := scanner.Bytes()
	res.dirs = make([]byte, len(dirs))
	for i, d := range dirs {
		if d == 'L' {
			res.dirs[i] = 0
		} else {
			res.dirs[i] = 1
		}
	}
	// empty lines
	scanner.Scan()
	// graph
	res.vertexIndex = make(map[string]int)
	res.childs = make([][2]int, 0, 700)
	for scanner.Scan() {
		line := scanner.Bytes()
		// example: "JKB = (RFF, TSX)"
		n_l_r := [3]string{string(line[0:3]), string(line[7:10]), string(line[12:15])}
		indices := [3]int{}
		for i := 0; i < 3; i++ {
			index, ok := res.vertexIndex[n_l_r[i]]
			if !ok {
				index = len(res.childs)
				res.childs = append(res.childs, [2]int{})
				res.vertexIndex[n_l_r[i]] = index
			}
			indices[i] = index
		}
		res.childs[indices[0]] = [2]int{indices[1], indices[2]}
	}
	return res
}

func solve1(in []byte) int {
	data := parseInput(in)
	pos := data.vertexIndex["AAA"]
	end := data.vertexIndex["ZZZ"]
	steps := 0
	for pos != end {
		pos = data.childs[pos][data.dirs[steps%len(data.dirs)]]
		steps++
	}
	return steps
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func solve2(in []byte) int {
	data := parseInput(in)
	starts := make([]int, 0, 100)
	ends := make([]int, 0, 100)
	for k, i := range data.vertexIndex {
		if k[2] == 'A' {
			starts = append(starts, i)
		} else if k[2] == 'Z' {
			ends = append(ends, i)
		}
	}

	// looking at the ghosts moves: each ghosts visit an end node with a fixed period
	// and each cycle starts at pos 0
	cycleLengths := make([]int, len(starts))
	for i, start := range starts {
		pos := start
		steps := 0
		for !slices.Contains(ends, pos) {
			pos = data.childs[pos][data.dirs[steps%len(data.dirs)]]
			steps++
		}
		cycleLengths[i] = steps
	}

	res := cycleLengths[0]
	for _, c := range cycleLengths[1:] {
		res = lcm(res, c)
	}

	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
