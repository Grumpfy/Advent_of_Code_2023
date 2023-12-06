package main

import (
	"fmt"
	"math"
	"time"
)

type Race struct {
	T int
	D int
}

var sampleInput = []Race{{7, 9}, {15, 40}, {30, 200}}
var input1 = []Race{{40, 277}, {82, 1338}, {91, 1349}, {66, 1063}}
var input2 = Race{40829166, 277133813491063}

func (r Race) ways() int {
	delta := math.Sqrt(float64(r.T*r.T - 4*r.D))
	t := float64(r.T)
	h := int(math.Floor((t+delta)/2.0))
	l := int(math.Ceil((t-delta)/2.0))
	if r.eval(h) == r.D {
		h--
	}
	if r.eval(l) == r.D {
		l++
	}
	return int(h-l) + 1
}

func (r Race) eval(t int) int {
	return t * (r.T - t)
}

func (r Race) lowerBound(start int, length int, val int) int {
	for length > 0 {
		steps := length / 2
		mid := start + steps

		if r.eval(mid) < val {
			start = mid + 1
			length -= steps + 1
		} else {
			length = steps
		}
	}
	return start
}

func (r Race) anotherWays() int {
	l := r.lowerBound(0, r.T/2, r.D+1)
	return 2*(r.T/2-l+1) - 1 + r.T%2
}

func solve1(in []Race) int {
	res := 1
	for _, r := range in {
		ways := r.ways()
		res *= ways
	}
	return res
}

func main() {
	start := time.Now()
	// res := solve1(input1)
	res := input2.ways()
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
