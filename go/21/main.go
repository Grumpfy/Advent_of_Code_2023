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

func solve1(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	startPos := bytes.IndexByte(in, 'S')
	positions := [2][]int{make([]int, 0, 1024), make([]int, 0, 1024)}
	current, next := 0, 1
	positions[next] = append(positions[next], startPos)
	for i := 0; i < 64; i++ {
		next, current = current, next
		positions[next] = positions[next][0:0]
		for _, pos := range positions[current] {
			if (pos%width) > 0 && in[pos-1] != '#' {
				positions[next] = append(positions[next], pos-1)
			}
			if (pos%width) < width-2 && in[pos+1] != '#' {
				positions[next] = append(positions[next], pos+1)
			}
			if pos >= width && in[pos-width] != '#' {
				positions[next] = append(positions[next], pos-width)
			}
			if pos+width < len(in)-1 && in[pos+width] != '#' {
				positions[next] = append(positions[next], pos+width)
			}
		}
		slices.Sort(positions[next])
		positions[next] = slices.Compact(positions[next])
	}
	return len(positions[next])
}

type Pos2D [2]int

func (p Pos2D) index(width, height, stride int) int {
	return stride*((p[1]%height+height)%height) + (p[0]%width+width)%width
}

func solve2(in []byte) int {
	width := bytes.IndexByte(in, '\n')
	stride := width + 1
	height := len(in) / stride
	startIndex := bytes.IndexByte(in, 'S')
	positions := [2][]Pos2D{make([]Pos2D, 0, 1024), make([]Pos2D, 0, 1024)}
	current, next := 0, 1
	dists := map[Pos2D]int{}
	positions[next] = append(positions[next], Pos2D{startIndex % stride, startIndex / stride})
	dists[Pos2D{startIndex % stride, startIndex / stride}] = 0
	for i := 0; i < 50; i++ {
		next, current = current, next
		positions[next] = positions[next][0:0]
		for _, pos := range positions[current] {
			neighbours := [4]Pos2D{
				{pos[0] - 1, pos[1]},
				{pos[0] + 1, pos[1]},
				{pos[0], pos[1] - 1},
				{pos[0], pos[1] + 1}}
			for _, n := range neighbours {
				if in[n.index(width, height, stride)] != '#' {
					if _, found := dists[n]; !found {
						dists[n] = i + 1
					}
					positions[next] = append(positions[next], n)
				}
			}
		}
		cmp := func(l, r Pos2D) int {
			if l[0] != r[0] {
				return l[0] - r[0]
			}
			return l[1] - r[1]
		}

		slices.SortFunc(positions[next], cmp)
		positions[next] = slices.Compact(positions[next])
		for j := -2 * height; j < 3*height; j++ {
			for i := -2 * width; i < 3*width; i++ {
				c := in[Pos2D{i, j}.index(width, height, stride)]
				if v, found := dists[Pos2D{i, j}]; found {
					fmt.Print(v % 2)
				} else {
					// fmt.Printf("\033[31m%s\033[0m", string(c))
					fmt.Print(string(c))
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
	return len(positions[next])
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
