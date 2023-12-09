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

func scanNextInt(in []byte) (int, []byte) {
	pos := 0
	end := len(in)
	for pos != end {
		if in[pos] == ' ' {
			pos++
		} else {
			break
		}
	}
	sign := 1
	if pos != end && in[pos] == '-' {
		sign = -1
		pos++
	}
	res := 0
	for pos != end {
		if in[pos] < '0' || in[pos] > '9' {
			break
		}
		res = res*10 + int(in[pos]-'0')
		pos++
	}
	return sign * res, in[pos:]
}

type historyScanner struct {
	scanner *bufio.Scanner
	history []int
}

func NewHistoryScanner(in []byte) *historyScanner {
	res := &historyScanner{}
	res.scanner = bufio.NewScanner(bytes.NewReader(in))
	res.history = make([]int, 0, 128)	
	return res
}

func (h *historyScanner) scan() bool {
	hasNext := h.scanner.Scan()
	if hasNext {
		h.history = h.history[:0]
		line := h.scanner.Bytes()
		for len(line) > 0 {
			var v int
			v, line = scanNextInt(line)
			h.history = append(h.history, v)
		}
	}
	return hasNext
}

func nextValue(history []int) int {
	end := len(history)
	for end > 1 {
		hasNonZero := false
		for i := 0; i < (end - 1); i++ {
			history[i] = history[i+1] - history[i]
			hasNonZero = hasNonZero || (history[i] != 0)
		}
		if !hasNonZero {
			break
		}
		end--
	}
	res := 0
	for _, v := range history {
		res += v
	}
	return res
}

func solve1(in []byte) int {
	scanner := NewHistoryScanner(in)
	res := 0
	for scanner.scan() {
		res += nextValue(scanner.history)
	}
	return res
}

func prevValue(history []int) int {
	end := len(history)
	for end > 1 {
		hasNonZero := false
		first := history[0]
		for i := 0; i < (end - 1); i++ {
			history[i] = history[i+1] - history[i]
			hasNonZero = hasNonZero || (history[i] != 0)
		}
		history[end-1] = first 
		if !hasNonZero {
			break
		}
		end--
	}
	res := 0
	for _, v := range history {
		res = v - res
	}
	return res
}

func solve2(in []byte) int {
	scanner := NewHistoryScanner(in)
	res := 0
	for scanner.scan() {
		res += prevValue(scanner.history)
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
