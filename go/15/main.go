package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

func hashCombine(seed int, c byte) int {
	seed += int(c)
	seed *= 17
	return seed % 256
}

func solve1(in []byte) int {
	res := 0
	hash := 0
	for _, c := range in {
		switch c {
		case ',', '\n':
			res += hash
			hash = 0
		default:
			hash = hashCombine(hash, c)
		}
	}
	return res
}

func parseLabel(in []byte, start int) ([]byte, int) {
	pos := start
	for end := len(in); pos < end; pos++ {
		if in[pos] < 'a' || in[pos] > 'z' {
			break
		}
	}
	return in[start:pos], pos
}

func parseInt(in []byte, start int) (int, int) {
	pos := start
	res := 0
	for end := len(in); pos < end; pos++ {
		if in[pos] < '0' || in[pos] > '9' {
			break
		}
		res = res*10 + int(in[pos]-'0')
	}
	return res, pos
}

func getBoxID(label []byte) int {
	hash := 0
	for _, c := range label {
		hash = hashCombine(hash, c)
	}
	return hash
}

func solve2(in []byte) int {
	type Lens struct {
		label []byte
		focal int
	}
	buckets := [256][]Lens{}
	for i := range buckets {
		buckets[i] = make([]Lens, 0, 128)
	}
	pos := 0
	for end := len(in); pos < end; {
		var label []byte
		label, pos = parseLabel(in, pos)
		boxID := getBoxID(label)
		switch in[pos] {
		case '-':
			for lens := range buckets[boxID] {
				if bytes.Compare(buckets[boxID][lens].label,  label) == 0 {
					buckets[boxID] = slices.Delete(buckets[boxID], lens, lens+1)
					break
				}
			}
			pos += 2
		case '=':
			focal := 0
			focal, pos = parseInt(in, pos+1)
			found := false
			for lens := range buckets[boxID] {
				if bytes.Compare(buckets[boxID][lens].label,  label) == 0 {
					buckets[boxID][lens].focal = focal
					found = true
					break
				}
			}
			if ! found {
				buckets[boxID] = append(buckets[boxID], Lens{label, focal})
			}
			pos++
		case '\n':
			pos++
		}
	}
	
	res := 0
	for boxID := range buckets {
		for lensID := range buckets[boxID] {
			res += (boxID + 1) * (lensID + 1) * buckets[boxID][lensID].focal
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
