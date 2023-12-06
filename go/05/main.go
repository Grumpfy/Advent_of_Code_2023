package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"math"
	"time"
)

//go:embed input.txt
var input []byte

var NotInRange = errors.New("not in range")

type Range struct {
	start  int
	length int
}

func (r Range) isEmpty() bool {
	return r.length <= 0
}

type Remap struct {
	destStart   int
	sourceStart int
	rangeLength int
}

func (r *Remap) apply(n int) (int, error) {
	d := n - r.sourceStart
	if d >= 0 && d < r.rangeLength {
		return r.destStart + d, nil
	}
	return 0, NotInRange
}

func (r *Remap) applyRange(in Range) (remap_ Range, rest_ Range) {
	d := in.start - r.sourceStart
	if d >= 0 && d < r.rangeLength {
		lenInRange := min(r.rangeLength - d, in.length)
		lenOutRange := in.length - lenInRange
		return Range{r.destStart + d, lenInRange}, Range{in.start + lenInRange, lenOutRange}
	}
	return Range{}, in
}

type Remaps []Remap

type PuzzleData struct {
	seeds []int
	maps  []Remaps
}

func NewPuzzleData() *PuzzleData {
	return &PuzzleData{
		seeds: make([]int, 0, 25),
		maps:  make([]Remaps, 0, 8),
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func parseInt(in []byte) (int, []byte) {
	pos, end, res := 0, len(in), 0
	for ; pos < end; pos++ {
		if isDigit(in[pos]) {
			res = 10*res + int(in[pos]-'0')
		} else {
			break
		}
	}
	return res, in[pos:]
}

func goToNextDigit(in []byte) []byte {
	pos, end := 0, len(in)
	for ; pos < end; pos++ {
		if isDigit(in[pos]) {
			return in[pos:]
		}
	}
	return in[pos:]
}

func ParseInput(in []byte) *PuzzleData {
	res := NewPuzzleData()
	scanner := bufio.NewScanner(bytes.NewReader(in))
	// pase seeds
	scanner.Scan()
	line := scanner.Bytes()
	seed := 0
	for len(line) > 0 {
		line = goToNextDigit(line)
		seed, line = parseInt(line)
		res.seeds = append(res.seeds, seed)
	}
	// parse range remaps
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			res.maps = append(res.maps, Remaps{})
			scanner.Scan() // skipp 'x-to-y map:'
		} else {
			r := Remap{}
			r.destStart, line = parseInt(line)
			line = goToNextDigit(line)
			r.sourceStart, line = parseInt(line)
			line = goToNextDigit(line)
			r.rangeLength, line = parseInt(line)
			end := len(res.maps) - 1
			res.maps[end] = append(res.maps[end], r)
		}
	}
	return res
}

func solve1(in []byte) int {
	data := ParseInput(in)
	minLocation := math.MaxInt
	for _, seed := range data.seeds {
		location := seed
		for m := range data.maps {
			for r := range data.maps[m] {
				dest, err := data.maps[m][r].apply(location)
				if err == nil {
					location = dest
					break
				}
			}
		}
		minLocation = min(minLocation, location)
	}
	return minLocation
}

func solve2(in []byte) int {
	data := ParseInput(in)
	minLocation := math.MaxInt
	for i, end := 0, len(data.seeds)/2; i < end; i++ {
		seedStart := data.seeds[2*i]
		seedEnd := seedStart + data.seeds[2*i+1]
		for seed := seedStart; seed < seedEnd; seed++ {
			location := seed
			for m := range data.maps {
				for r := range data.maps[m] {
					dest, err := data.maps[m][r].apply(location)
					if err == nil {
						location = dest
						break
					}
				}
			}
			minLocation = min(minLocation, location)
		}
	}
	return minLocation
}

func solve2b(in []byte) int {
	data := ParseInput(in)
	minLocation := math.MaxInt
	for i, end := 0, len(data.seeds)/2; i < end; i++ {
		inRanges := make([]Range, 0, 128)
		outRanges := make([]Range, 0, 128)
		inRanges = append(inRanges, Range{data.seeds[2*i], data.seeds[2*i+1]})
		for m := range data.maps {
			for len(inRanges) != 0 {
				inRange := inRanges[len(inRanges)-1]
				inRanges = inRanges[:len(inRanges)-1]
				consumed := false
				for r := range data.maps[m] {
					outRange, rest := data.maps[m][r].applyRange(inRange)
					if !outRange.isEmpty() {
						outRanges = append(outRanges, outRange)
						if !rest.isEmpty() {
							inRanges = append(inRanges, rest)
						}
						consumed = true
						break
					}
				}
				if !consumed {
					outRanges = append(outRanges, inRange)
				}
			}
			inRanges, outRanges = outRanges, inRanges
		}
		for _, r := range inRanges {
			minLocation = min(minLocation, r.start)
		}
	}
	return minLocation
}

func main() {
	start := time.Now()
	res := solve2b(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
