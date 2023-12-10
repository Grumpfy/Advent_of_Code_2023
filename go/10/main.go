package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var input []byte

type Dir [2]int

func next(pipesMap []byte, width int, pos int, dir Dir) (int, Dir) {
	switch pipesMap[pos] {
	case '|', '-':
	case 'L': // 0,1 => 1,0 | -1,0 => 0,-1
		dir = Dir{dir[1], dir[0]}
	case 'J': // 0,1 => -1,0 | 1,0 => 0,-1
		dir = Dir{-1 * dir[1], -1 * dir[0]}
	case '7': // 1,0 => 1,0 | -1,0 => 0,-1
		dir = Dir{dir[1], dir[0]}
	case 'F': // -1,0 => 0,1 | 0,-1 => 1,0
		dir = Dir{-1 * dir[1], -1 * dir[0]}
	case 'S':
		if bytes.IndexByte([]byte{'|', 'F', '7'}, pipesMap[pos-width]) >= 0 {
			dir = Dir{0, -1}
		} else if bytes.IndexByte([]byte{'-', 'L', 'F'}, pipesMap[pos-1]) >= 0 {
			dir = Dir{-1, 0}
		} else if bytes.IndexByte([]byte{'|', 'L', 'J'}, pipesMap[pos+width]) >= 0 {
			dir = Dir{0, 1}
		} else if bytes.IndexByte([]byte{'-', '7', 'J'}, pipesMap[pos+1]) >= 0 {
			dir = Dir{1, 0}
		} else {
			panic("")
		}
	default:
		panic("")
	}
	return pos + dir[0] + width*dir[1], dir
}

func solve1(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	pos := bytes.IndexByte(in, 'S')
	dir := Dir{}
	steps := 1
	pos, dir = next(in, width, pos, dir)
	for in[pos] != 'S' {
		pos, dir = next(in, width, pos, dir)
		steps++
	}
	return steps / 2
}

func solve2(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	pos := bytes.IndexByte(in, 'S')
	dir := Dir{}
	cleanMap := make([]byte, len(in))
	cleanMap[pos] = '7' // just looked at the input
	pos, dir = next(in, width, pos, dir)
	for in[pos] != 'S' {
		cleanMap[pos] = in[pos]
		pos, dir = next(in, width, pos, dir)
	}
	i := 0
	inside := false
	res := 0
	for i < len(cleanMap) {
		for _, c := range cleanMap[i : i+width] {
			switch c {
			case 0:
				if inside {
					res += 1
					// fmt.Print(".")
				} else {
					// fmt.Print(" ")
				}
			case '-', '7', 'F':
				// fmt.Print(string(c))
			case '|', 'L', 'J':
				inside = !inside
				// fmt.Print(string(c))
			default:
				panic("")
			}
		}
		// fmt.Println()
		i += width
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
