package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"time"
)

//go:embed input.txt
var input []byte

const (
	Right = uint8(0)
	Down  = uint8(1)
	Left  = uint8(2)
	Up    = uint8(3)
)

type Ray struct {
	pos int
	dir uint8
}

func advance(ray Ray, width int, end int) (Ray, error) {
	switch ray.dir {
	case Right:
		if (ray.pos % width) == width-2 {
			return Ray{}, errors.New("invalid pos")
		}
		ray.pos += 1
	case Down:
		ray.pos += width
		if ray.pos >= (end - 1) {
			return Ray{}, errors.New("invalid pos")
		}
	case Left:
		if (ray.pos % width) == 0 {
			return Ray{}, errors.New("invalid pos")
		}
		ray.pos -= 1
	case Up:
		ray.pos -= width
		if ray.pos < 0 {
			return Ray{}, errors.New("invalid pos")
		}
	}
	return ray, nil
}

func solve(in []byte, start Ray) int {
	width := bytes.IndexByte(in, '\n') + 1
	edges := make([]uint8, len(in))
	rays := make([]Ray, 0, 1024)
	rays = append(rays, start)

	for len(rays) != 0 {
		ray := rays[len(rays)-1]
		rays = rays[:len(rays)-1]
		if edges[ray.pos]&(1<<ray.dir) == 0 {
			edges[ray.pos] |= (1 << ray.dir)
			switch in[ray.pos] {
			case '.':
				if newRay, err := advance(ray, width, len(in)); err == nil {
					rays = append(rays, newRay)
				}
			case '-':
				if ray.dir&1 == 0 {
					if newRay, err := advance(ray, width, len(in)); err == nil {
						rays = append(rays, newRay)
					}
				} else {
					rays = append(rays, Ray{ray.pos, Left}, Ray{ray.pos, Right})
				}
			case '|':
				if ray.dir&1 != 0 {
					if newRay, err := advance(ray, width, len(in)); err == nil {
						rays = append(rays, newRay)
					}
				} else {
					rays = append(rays, Ray{ray.pos, Up}, Ray{ray.pos, Down})
				}
			case '\\':
				// switch ray.dir {
				// case Right: // 00 => 01
				// 	ray.dir = Down
				// case Down:  // 01 => 00
				// 	ray.dir = Right
				// case Left:  // 10 => 11
				// 	ray.dir = Up
				// case Up:    // 11 => 10
				// 	ray.dir = Left
				// }
				ray.dir ^= 1
				if newRay, err := advance(ray, width, len(in)); err == nil {
					rays = append(rays, newRay)
				}
			case '/':
				// switch ray.dir {
				// case Right: // 00 => 11
				// 	ray.dir = Up
				// case Down:  // 01 => 10
				// 	ray.dir = Left
				// case Left:  // 10 => 01
				// 	ray.dir = Down
				// case Up:    // 11 => 00
				// 	ray.dir = Right
				// }
				ray.dir ^= 3
				if newRay, err := advance(ray, width, len(in)); err == nil {
					rays = append(rays, newRay)
				}
			}
		}
	}

	res := 0
	for _, c := range edges {
		if c != 0 {
			res++
		}
	}
	return res
}

func solve1(in []byte) int {
	return solve(in, Ray{0, Right})
}

func solve2(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	height := (len(in) + width-1) % width
	res := 0
	for i := 0; i < width-1; i++ {
		res = max(res, solve(in, Ray{i, Down}))
		res = max(res, solve(in, Ray{i + width * (height - 1), Up}))
	}
	for i := 0; i < height; i++ {
		res = max(res, solve(in, Ray{i * width, Right}))
		res = max(res, solve(in, Ray{(i+1) * width - 2, Left}))
	}
	return res
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
