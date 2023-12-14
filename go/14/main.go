package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/spaolacci/murmur3"
	"time"
)

//go:embed input.txt
var input []byte

func load(in []byte, height int) int {
	weight := height
	res := 0
	for _, c := range in {
		if c == 'O' {
			res += weight
		} else if c == '\n' {
			weight--
		}
	}
	return res
}

func solve1(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	for i, end := width, len(in); i < end; i++ {
		if in[i] == 'O' {
			j := i - width
			for ; j >= 0; j -= width {
				if in[j] != '.' {
					break
				}
			}
			j += width
			in[i], in[j] = in[j], in[i]
		}
	}
	height := (len(in) + width - 1) / width
	return load(in, height)
}

func rollLine(in []byte, begin int, end int, stride int) {
	for i := begin + stride; i != end; i += stride {
		if in[i] == 'O' {
			j := i - stride
			for ; j != begin; j -= stride {
				if in[j] != '.' {
					break
				}
			}
			if in[j] != '.' {
				j += stride
			}
			in[i], in[j] = in[j], in[i]
		}
	}
}

func cycle(in []byte, width int, height int) {
	// north
	for i := 0; i < width-1; i++ {
		rollLine(in, i, i+height*width, width)
	}
	// west
	for i := 0; i < height; i++ {
		rollLine(in, i*width, (i+1)*width, 1)
	}
	// south
	for i := 0; i < width-1; i++ {
		rollLine(in, i+(height-1)*width, i-width, -width)
	}
	// east
	for i := 0; i < height; i++ {
		rollLine(in, (i+1)*width-1, i*width-1, -1)
	}
}

func solve2(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	height := (len(in) + width - 1) / width

	knownConfigurations := make(map[uint64]int)
	steps := 0
	currentConf := murmur3.Sum64(in)
	cycleStart, cycleLenght := 0, 0
	for {
		if dejaVu, ok := knownConfigurations[currentConf]; !ok {
			knownConfigurations[currentConf] = steps
			cycle(in, width, height)
			steps++
			currentConf = murmur3.Sum64(in)
		} else {
			cycleStart = dejaVu
			cycleLenght = steps - dejaVu
			break
		}
	}

	remainingSteps := (1000000000 - cycleStart) % cycleLenght
	for i := 0; i < remainingSteps; i++ {
		cycle(in, width, height)
	}

	return load(in, height)
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
