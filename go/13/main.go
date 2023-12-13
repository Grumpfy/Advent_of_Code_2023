package main

import (
	_ "embed"
	"fmt"
	"math/bits"
	"time"
)

//go:embed input.txt
var input []byte

func symmetryPos(l []uint) int {
	for i, end := 1, len(l); i < end; i++ {
		left, right := i-1, i
		match := true
		for left >= 0 && right < end {
			if l[left] != l[right] {
				match = false
				break
			}
			left--
			right++
		}
		if match {
			return i
		}
	}
	return -1
}

func summary(lines, columns []uint) int {
	if xSym := symmetryPos(columns); xSym != -1 {
		return xSym
	}
	return symmetryPos(lines) * 100
}

func solve1(in []byte) int {
	lines := make([]uint, 0, 128)
	columns := make([]uint, 0, 128)
	colID := 0
	lineValue := uint(0)
	res := 0
	for _, c := range in {
		switch c {
		case '.':
			lineValue = lineValue << 1
			if len(columns) <= colID {
				columns = append(columns, 0)
			}
			columns[colID] = columns[colID] << 1
			colID++
		case '#':
			lineValue = (lineValue << 1) | 1
			if len(columns) <= colID {
				columns = append(columns, 0)
			}
			columns[colID] = (columns[colID] << 1) | 1
			colID++
		case '\n':
			if colID == 0 {
				// end of pattern
				res += summary(lines, columns)
				lines = lines[:0]
				columns = columns[:0]
			} else {
				// end of line
				lines = append(lines, lineValue)
				colID = 0
				lineValue = 0
			}
		}
	}
	res += summary(lines, columns)
	return res
}

func symmetrySmudgePos(l []uint) int {
	for i, end := 1, len(l); i < end; i++ {
		left, right := i-1, i
		match := true
		smudge := false
		for left >= 0 && right < end {
			if l[left] != l[right] {
				if smudge || bits.OnesCount(l[left]^l[right]) != 1 {
					match = false
					break
				}
				smudge = true
			}
			left--
			right++
		}
		if match && smudge {
			return i
		}
	}
	return -1
}

func summarySmudge(lines, columns []uint) int {
	if xSym := symmetrySmudgePos(columns); xSym != -1 {
		return xSym
	}
	return symmetrySmudgePos(lines) * 100
}

func solve2(in []byte) int {
	lines := make([]uint, 0, 128)
	columns := make([]uint, 0, 128)
	colID := 0
	lineValue := uint(0)
	res := 0
	for _, c := range in {
		switch c {
		case '.':
			lineValue = lineValue << 1
			if len(columns) <= colID {
				columns = append(columns, 0)
			}
			columns[colID] = columns[colID] << 1
			colID++
		case '#':
			lineValue = (lineValue << 1) | 1
			if len(columns) <= colID {
				columns = append(columns, 0)
			}
			columns[colID] = (columns[colID] << 1) | 1
			colID++
		case '\n':
			if colID == 0 {
				// end of pattern
				res += summarySmudge(lines, columns)
				lines = lines[:0]
				columns = columns[:0]
			} else {
				// end of line
				lines = append(lines, lineValue)
				colID = 0
				lineValue = 0
			}
		}
	}
	res += summarySmudge(lines, columns)
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
