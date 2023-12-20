package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

type Pos2D [2]int

func (p *Pos2D) Add(offset Pos2D) {
	p[0] += offset[0]
	p[1] += offset[1]
}

func (p *Pos2D) Madd(offset Pos2D, l int) {
	p[0] += l * offset[0]
	p[1] += l * offset[1]
}

func PieceWiseMin(l Pos2D, r Pos2D) Pos2D {
	return Pos2D{
		min(l[0], r[0]),
		min(l[1], r[1]),
	}
}

func PieceWiseMax(l Pos2D, r Pos2D) Pos2D {
	return Pos2D{
		max(l[0], r[0]),
		max(l[1], r[1]),
	}
}

type Dig struct {
	trenchDir Pos2D
	trenchLen int
}

func parseInt(in []byte, startPos int) (int, int) {
	pos := startPos
	res := 0
	for ; pos < len(in); pos++ {
		if in[pos] < '0' || in[pos] > '9' {
			break
		}
		res = res*10 + int(in[pos]-'0')
	}
	return res, pos
}

func parseInput(in []byte) []Dig {
	res := make([]Dig, 0, 1024)
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Bytes()
		offsetLen, _ := parseInt(line, 2)
		var offset Pos2D
		switch line[0] {
		case 'L':
			offset = Pos2D{-1, 0}
		case 'R':
			offset = Pos2D{1, 0}
		case 'U':
			offset = Pos2D{0, -1}
		case 'D':
			offset = Pos2D{0, 1}
		}
		res = append(res, Dig{offset, offsetLen})
	}
	return res
}

func solve1(in []byte) int {
	digPlan := parseInput(in)
	bbMin, bbMax := Pos2D{math.MaxInt, math.MaxInt}, Pos2D{math.MinInt, math.MinInt}
	{
		currentPos := Pos2D{}
		for _, d := range digPlan {
			currentPos.Madd(d.trenchDir, d.trenchLen)
			bbMax = PieceWiseMax(bbMax, currentPos)
			bbMin = PieceWiseMin(bbMin, currentPos)
		}
	}
	dim := Pos2D{
		bbMax[0] - bbMin[0] + 1,
		bbMax[1] - bbMin[1] + 1,
	}
	digMap := make([]byte, dim[0]*dim[1])
	{
		currentPos := Pos2D{-bbMin[0], -bbMin[1]}
		digMap[currentPos[0]+dim[0]*currentPos[1]] = '#'
		for _, d := range digPlan {
			for i := 0; i < d.trenchLen; i++ {
				currentPos.Add(d.trenchDir)
				digMap[currentPos[0]+dim[0]*currentPos[1]] = '#'
			}
		}
	}

	frontier := make([]int, 1, 1024)
	// find interior point
	for j := 0; j < dim[1]; j++ {
		line := digMap[j*dim[0] : (j+1)*dim[0]-1]
		idx := bytes.IndexByte(line, '#')
		if idx != -1 && line[idx+1] == 0 {
			frontier[0] = j*dim[0] + idx + 1
			break
		}
	}
	// fill
	for len(frontier) > 0 {
		pos := frontier[len(frontier)-1]
		frontier = frontier[0 : len(frontier)-1]
		if digMap[pos] == '#' {
			continue
		}
		digMap[pos] = '#'
		frontier = append(frontier, pos-1, pos+1, pos-dim[0], pos+dim[0])
	}
	// count
	res := 0
	for _, c := range digMap {
		if c == '#' {
			res++
		}
	}

	return res
}

func parseInput2(in []byte) []Dig {
	res := make([]Dig, 0, 1024)
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Bytes()
		startPos := bytes.IndexByte(line, '#') + 1
		offsetLen := 0
		for i := 0; i < 5; i++ {
			c := line[startPos+i]
			if c >= '0' && c <= '9' {
				offsetLen = 16*offsetLen + int(c-'0')
			} else {
				offsetLen = 16*offsetLen + int(c-'a') + 10 
			}
		}
		var offset Pos2D
		switch line[startPos+5] {
		case '0':
			offset = Pos2D{1, 0}
		case '1':
			offset = Pos2D{0, 1}
		case '2':
			offset = Pos2D{-1, 0}
		case '3':
			offset = Pos2D{0, -1}
		}
		res = append(res, Dig{offset, offsetLen})
	}
	return res
}

func solve2(in []byte) int {
	digPlan := parseInput2(in)
	xAxis := make([]int, 1, 1024) // [0]
	yAxis := make([]int, 1, 1024) // [0]
	{
		currentPos := Pos2D{}
		for _, d := range digPlan {
			currentPos.Madd(d.trenchDir, d.trenchLen)
			xAxis = append(xAxis, currentPos[0], currentPos[0]+1)
			yAxis = append(yAxis, currentPos[1], currentPos[1]+1)
		}
		slices.Sort(xAxis)
		xAxis = slices.Compact(xAxis)
		slices.Sort(yAxis)
		yAxis = slices.Compact(yAxis)
		// fmt.Println(xAxis, yAxis)
	}
	dim := Pos2D{
		len(xAxis),
		len(yAxis),
	}
	digMap := make([]byte, dim[0]*dim[1])
	{
		currentPos := Pos2D{}
		i, _ := slices.BinarySearch(xAxis, currentPos[0])
		j, _ := slices.BinarySearch(yAxis, currentPos[1])
		digMap[i+dim[0]*j] = '#'
		// fmt.Println(currentPos, i,j)
		for _, d := range digPlan {
			currentPos.Madd(d.trenchDir, d.trenchLen)
			if d.trenchDir[0] == 0 {
				endJ, _ := slices.BinarySearch(yAxis, currentPos[1])
				for k, end := min(j, endJ), max(j, endJ); k <= end; k++ {
					digMap[i+dim[0]*k] = '#'
				}
				j = endJ
			} else {
				endI, _ := slices.BinarySearch(xAxis, currentPos[0])
				for k, end := min(i, endI), max(i, endI); k <= end; k++ {
					digMap[k+dim[0]*j] = '#'
				}
				i = endI
			}
			// fmt.Println(currentPos, i,j)
		}
	}

	// for j := 0; j < dim[1]-1; j++ {
	// 	for i := 0; i < dim[0]-1; i++ {
	// 		if digMap[i+j*dim[0]] == '#' {
	// 			fmt.Print(string('#'))
	// 		} else {
	// 			fmt.Print(string(' '))
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	frontier := make([]int, 1, 1024)
	// find interior point
	for j := 0; j < dim[1]; j++ {
		line := digMap[j*dim[0] : (j+1)*dim[0]-1]
		idx := bytes.IndexByte(line, '#')
		if idx != -1 && line[idx+1] == 0 {
			frontier[0] = j*dim[0] + idx + 1
			break
		}
	}
	// fill
	for len(frontier) > 0 {
		pos := frontier[len(frontier)-1]
		frontier = frontier[0 : len(frontier)-1]
		if digMap[pos] == '#' {
			continue
		}
		digMap[pos] = '#'
		frontier = append(frontier, pos-1, pos+1, pos-dim[0], pos+dim[0])
	}
	// count
	res := 0
	for j := 0; j < dim[1]-1; j++ {
		for i := 0; i < dim[0]-1; i++ {
			if digMap[i+j*dim[0]] == '#' {
				res += (xAxis[i+1] - xAxis[i]) * (yAxis[j+1] - yAxis[j])
				// fmt.Print(string('#'))
			} else {
				// fmt.Print(string(' '))
			}
		}
		// fmt.Println()
	}

	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
