package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"time"
)

//go:embed input.txt
var input []byte

type Pos3D [3]int
type Brick [2]Pos3D

func parseInput(in []byte) []Brick {
	res := make([]Brick, 0, 1394)
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		// 4,6,217~4,8,217
		line := scanner.Bytes()
		var brick Brick
		for i, pos := range bytes.Split(line, []byte{'~'}) {
			for j, coord := range bytes.Split(pos, []byte{','}) {
				brick[i][j], _ = strconv.Atoi(string(coord))
			}
		}
		res = append(res, brick)
	}
	return res
}

type ZBuff [10][10]int

func solve1(in []byte) int {
	bricks := parseInput(in)
	z := ZBuff{}
	id := ZBuff{}
	for j := 0; j < 10; j++ {
		for i := 0; i < 10; i++ {
			id[i][j] = -1
		}
	}

	slices.SortFunc(bricks, func(l, r Brick) int { return l[0][2] - r[0][2] })
	unsafeBricks := make([]int, 0, 2048)
	for i, b := range bricks {
		zmax := 0
		idsAtMax := make([]int, 0, 100)
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				if z[x][y] > zmax {
					zmax = z[x][y]
					idsAtMax = idsAtMax[0:0]
					idsAtMax = append(idsAtMax, id[x][y])
				} else if z[x][y] == zmax && !slices.Contains(idsAtMax, id[x][y]) {
					idsAtMax = append(idsAtMax, id[x][y])
				}
				zmax = max(zmax, z[x][y])
			}
		}
		if len(idsAtMax) == 1 && idsAtMax[0] != -1 {
			unsafeBricks = append(unsafeBricks, idsAtMax[0])
		}
		bricks[i][0][2] = zmax
		bricks[i][1][2] += zmax - b[0][2]
		topZ := bricks[i][1][2] + 1
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				z[x][y] = topZ
				id[x][y] = i
			}
		}

	}

	slices.Sort(unsafeBricks)
	unsafeBricks = slices.Compact(unsafeBricks)

	return len(bricks) - len(unsafeBricks)
}

func moveDown(bricks []Brick) int {
	z := ZBuff{}
	moves := 0
	for i, b := range bricks {
		zmax := 0
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				zmax = max(zmax, z[x][y])
			}
		}
		if b[0][2] != zmax {
			bricks[i][0][2] = zmax
			bricks[i][1][2] += zmax - b[0][2]
			moves++
		}
		topZ := bricks[i][1][2] + 1
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				z[x][y] = topZ
			}
		}
	}
	return moves
}

func solve2(in []byte) int {
	bricks := parseInput(in)
	slices.SortFunc(bricks, func(l, r Brick) int { return l[0][2] - r[0][2] })

	z := ZBuff{}
	id := ZBuff{}
	for j := 0; j < 10; j++ {
		for i := 0; i < 10; i++ {
			id[i][j] = -1
		}
	}
	unsafeBricks := make([]int, 0, 2048)
	for i, b := range bricks {
		zmax := 0
		idsAtMax := make([]int, 0, 100)
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				if z[x][y] > zmax {
					zmax = z[x][y]
					idsAtMax = idsAtMax[0:0]
					idsAtMax = append(idsAtMax, id[x][y])
				} else if z[x][y] == zmax && !slices.Contains(idsAtMax, id[x][y]) {
					idsAtMax = append(idsAtMax, id[x][y])
				}
				zmax = max(zmax, z[x][y])
			}
		}
		if len(idsAtMax) == 1 && idsAtMax[0] != -1 {
			unsafeBricks = append(unsafeBricks, idsAtMax[0])
		}
		bricks[i][0][2] = zmax
		bricks[i][1][2] += zmax - b[0][2]
		topZ := bricks[i][1][2] + 1
		for x := b[0][0]; x <= b[1][0]; x++ {
			for y := b[0][1]; y <= b[1][1]; y++ {
				z[x][y] = topZ
				id[x][y] = i
			}
		}

	}

	slices.Sort(unsafeBricks)
	unsafeBricks = slices.Compact(unsafeBricks)

	res := 0
	for _, b := range unsafeBricks {
		remaining := make([]Brick, 0, len(bricks)-1)
		remaining = append(remaining, bricks[0:b]...)
		remaining = append(remaining, bricks[b+1:]...)
		res += moveDown(remaining)
	}

	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
