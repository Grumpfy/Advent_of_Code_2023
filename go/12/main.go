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

func parseInt(in []byte) (int, []byte) {
	pos := 0
	end := len(in)
	res := 0
	for ; pos < end; pos++ {
		if in[pos] < '0' || in[pos] > '9' {
			break
		}
		res = res*10 + int(in[pos]-'0')
	}
	return res, in[pos:]
}

func parseLine(line []byte) ([]byte, []byte) {
	springEnd := bytes.IndexByte(line, ' ')
	spring := line[0:springEnd]
	pattern := make([]byte, 1, 128)
	pattern[0] = '*'
	line = line[springEnd+1:]
	first := true
	for len(line) > 0 {
		if !first {
			pattern = append(pattern, '.')
			pattern = append(pattern, '*')
			line = line[1:] // ','
		}
		first = false
		n := 0
		n, line = parseInt(line)
		for i := 0; i < n; i++ {
			pattern = append(pattern, '#')
		}
	}
	pattern = append(pattern, '*')
	return spring, pattern
}

func nbArrangements(spring, pattern []byte) int {
	lspring := len(spring)
	lpattern := len(pattern)
	w := lspring + 1
	mem := make([]int, (lpattern+1)*w)
	mem[(lpattern-1)*w+lspring] = 1
	for p := lpattern - 1; p >= 0; p-- {
		for s := lspring - 1; s >= 0; s-- {
			switch pattern[p] {
			case '*':
				switch spring[s] {
				case '?':
					mem[p*w+s] = mem[p*w+s+1] + mem[(p+1)*w+s]
				case '#':
					mem[p*w+s] = mem[(p+1)*w+s]
				case '.':
					mem[p*w+s] = mem[p*w+s+1]
				}
			case '#':
				switch spring[s] {
				case '?':
					mem[p*w+s] = mem[(p+1)*w+s+1]
				case '#':
					mem[p*w+s] = mem[(p+1)*w+s+1]
				case '.':
					mem[p*w+s] = 0
				}
			case '.':
				switch spring[s] {
				case '?':
					mem[p*w+s] = mem[(p+1)*w+s+1]
				case '#':
					mem[p*w+s] = 0
				case '.':
					mem[p*w+s] = mem[(p+1)*w+s+1]
				}
			}
		}
	}
	return mem[0]
}

func solve1(in []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(in))
	res := 0
	for scanner.Scan() {
		spring, pattern := parseLine(scanner.Bytes())
		res += nbArrangements(spring, pattern)
	}
	return res
}

func expand(spring, pattern []byte) ([]byte, []byte) {
	// spring
	expandedSpring := make([]byte, len(spring)*5+4)
	first := true
	pos := 0
	for pos < len(expandedSpring) {
		if !first {
			expandedSpring[pos] = '?'
			pos++
		}
		first = false
		copy(expandedSpring[pos:], spring)
		pos += len(spring)
	}
	// pattern
	trimmedPattern := pattern[1 : len(pattern)-1]
	expandedPattern := make([]byte, len(trimmedPattern)*5+10)
	first = true
	pos = 1
	expandedPattern[0] = '*'
	expandedPattern[len(expandedPattern)-1] = '*'
	for pos < len(expandedPattern)-1 {
		if !first {
			expandedPattern[pos] = '.'
			expandedPattern[pos+1] = '*'
			pos += 2
		}
		first = false
		copy(expandedPattern[pos:], trimmedPattern)
		pos += len(trimmedPattern)
	}
	return expandedSpring, expandedPattern
}

func solve2(in []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(in))
	res := 0
	for scanner.Scan() {
		spring, pattern := parseLine(scanner.Bytes())
		spring, pattern = expand(spring, pattern)
		res += nbArrangements(spring, pattern)
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
