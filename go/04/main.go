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

func skipSpaces(in []byte) []byte {
	for i, c := range in {
		if c != ' ' {
			return in[i:]
		}
	}
	return in[:0]
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func parseInt(in []byte) (int, []byte) {
	res := 0
	for i, c := range in {
		if isDigit(c) {
			res = 10*res + int(c-'0')
		} else {
			return res, in[i:]
		}
	}
	return res, in[:0]
}

func solve1(in []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(in))
	winningNumbers := make([]int, 0, 10)
	res := 0
	for scanner.Scan() {
		winningNumbers = winningNumbers[:0]
		line := scanner.Bytes()[9:] // skip 'Card XXX'
		for line[0] != '|' {
			line = skipSpaces(line)
			num := 0
			num, line = parseInt(line)
			winningNumbers = append(winningNumbers, num)
		}
		matches := 0
		line = line[1:] // skip '|'
		for len(line) != 0 {
			line = skipSpaces(line)
			num := 0
			num, line = parseInt(line)
			if slices.Contains(winningNumbers, num) {
				matches++
			}
		}
		if matches > 0 {
			res += 1 << (matches - 1)
		}
	}
	return res
}

func solve2(in []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(in))
	winningNumbers := make([]int, 0, 10)
	cards := make([]int, 202) // from input file
	for i := range cards {
		cards[i] = 1
	}
	cardID := 0
	for scanner.Scan() {
		winningNumbers = winningNumbers[:0]
		line := scanner.Bytes()[9:] // skip 'Card XXX'
		for line[0] != '|' {
			line = skipSpaces(line)
			num := 0
			num, line = parseInt(line)
			winningNumbers = append(winningNumbers, num)
		}
		matches := 0
		line = line[1:] // skip '|'
		for len(line) != 0 {
			line = skipSpaces(line)
			num := 0
			num, line = parseInt(line)
			if slices.Contains(winningNumbers, num) {
				matches++
			}
		}
		for i := 1; i <= matches; i++ {
			cards[cardID+i] += cards[cardID]
		}
		cardID++
	}
	res := 0
	for _, count := range cards {
		res += count
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
