package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var input []byte

func parseInt(b []byte) (int, []byte) {
	res := 0
	for i, c := range b {
		if c >= '0' && c <= '9' {
			res = 10*res + int(c-'0')
		} else {
			return res, b[i:]
		}
	}
	return 0, b
}

func solve1(in []byte) int {
	res := 0
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		// ex: 'Game 1: 7 green, 4 blue, 3 red; 4 blue, 10 red, 1 green; 1 blue, 9 red'
		rest := scanner.Bytes()
		rest = rest[5:] // 'Game '
		gameID := 0
		gameID, rest = parseInt(rest)
		isValid := true
		for isValid && len(rest) != 0 {
			rest = rest[2:] // ': ' or ', '
			nbCubes := 0
			nbCubes, rest = parseInt(rest)
			switch rest[1] {
			case 'r':
				if nbCubes > 12 {
					isValid = false
				}
				rest = rest[4:] // ' red'
			case 'g':
				if nbCubes > 13 {
					isValid = false
				}
				rest = rest[6:] // ' green'
			case 'b':
				if nbCubes > 14 {
					isValid = false
				}
				rest = rest[5:] // ' blue'
			}
		}
		if isValid {
			res += gameID
		}
	}
	return res
}

func solve2(in []byte) int {
	res := 0
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		// ex: 'Game 1: 7 green, 4 blue, 3 red; 4 blue, 10 red, 1 green; 1 blue, 9 red'
		rest := scanner.Bytes()
		rest = rest[5:] // 'Game '
		_, rest = parseInt(rest)
		maxR, maxG, maxB := 0, 0, 0
		for len(rest) != 0 {
			rest = rest[2:] // ': ' or ', '
			nbCubes := 0
			nbCubes, rest = parseInt(rest)
			switch rest[1] {
			case 'r':
				maxR = max(maxR, nbCubes)
				rest = rest[4:] // ' red'
			case 'g':
				maxG = max(maxG, nbCubes)
				rest = rest[6:] // ' green'
			case 'b':
				maxB = max(maxB, nbCubes)
				rest = rest[5:] // ' blue'
			}
		}
		res += maxR * maxG * maxB
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
