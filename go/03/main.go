package main

import (
	_ "embed"
	"fmt"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func parseNumber(in []byte) int {
	res, pos, end := 0, 0, len(in)
	for pos < end && isDigit(in[pos]) {
		res = 10*res + int(in[pos]-'0')
		pos++
	}
	return res
}

func solve1(in []byte) int {
	width := slices.Index(in, '\n') + 1
	offsets := []int{-width - 1, -width, -width + 1, -1, 1, width - 1, width, width + 1}
	end := len(in)
	parts := make([]int, 0, 1000)
	for i, c := range in {
		if !isDigit(c) && c != '.' && c != '\n' {
			// symbol
			for _, offset := range offsets {
				pos := i + offset
				if pos >= 0 && pos < end && isDigit(in[pos]) {
					for (pos-1) >= 0 && isDigit(in[pos-1]) {
						pos--
					}
					parts = append(parts, pos)
				}
			}
		}
	}
	slices.Sort(parts)
	parts = slices.Compact(parts)
	res := 0
	for _, part := range parts {
		res += parseNumber(in[part:])
	}
	return res
}

func solve2(in []byte) int {
	width := slices.Index(in, '\n') + 1
	offsets := []int{-width - 1, -width, -width + 1, -1, 1, width - 1, width, width + 1}
	end := len(in)
	parts := make([]int, 0, 1000)
	res := 0
	for i, c := range in {
		if c == '*' {
			parts = parts[:0]
			for _, offset := range offsets {
				pos := i + offset
				if pos >= 0 && pos < end && isDigit(in[pos]) {
					for (pos-1) >= 0 && isDigit(in[pos-1]) {
						pos--
					}
					parts = append(parts, pos)
				}
			}
			parts = slices.Compact(parts)
			if len(parts) == 2 {
				ratio := 1
				for _, part := range parts {
					ratio *= parseNumber(in[part:])
				}
				res += ratio
			}
		}
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
