package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"time"
)

//go:embed input.txt
var input []byte

type Game struct {
	hand       [5]byte
	handWeight int
	bid        int
}

func computeWeight(hand [5]byte) int {
	histogram := make([]int, 14)
	for i := 0; i < 5; i++ {
		histogram[hand[i]]++
	}
	slices.Sort(histogram[1:])
	return (histogram[13]+histogram[0])*10 + histogram[12]
}

func parseGames(in []byte, joker byte) []Game {
	res := make([]Game, 0, 1000)
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Bytes()
		newGame := Game{}
		for i := 0; i < 5; i++ {
			if line[i] == joker {
				newGame.hand[i] = 0
				continue
			}
			// A, K, Q, J, T, 9, 8, 7, 6, 5, 4, 3, or 2.
			switch line[i] {
			case 'A':
				newGame.hand[i] = 13
			case 'K':
				newGame.hand[i] = 12
			case 'Q':
				newGame.hand[i] = 11
			case 'J':
				newGame.hand[i] = 10
			case 'T':
				newGame.hand[i] = 9
			default:
				newGame.hand[i] = line[i] - '1'
			}
		}
		newGame.handWeight = computeWeight(newGame.hand)
		for _, c := range line[6:] {
			newGame.bid = newGame.bid*10 + int(c-'0')
		}
		res = append(res, newGame)
	}
	return res
}

func solve(in []byte, jocker byte) int {
	games := parseGames(in, jocker)
	slices.SortFunc(games, func(a, b Game) int {
		if a.handWeight != b.handWeight {
			return a.handWeight - b.handWeight
		}
		for i := 0; i < 5; i++ {
			if a.hand[i] != b.hand[i] {
				return int(a.hand[i]) - int(b.hand[i])
			}
		}
		return 0
	})

	res := 0
	for i, game := range games {
		res += (i+1) * game.bid
	}

	return res
}

func solve1(in []byte) int {
	return solve(in, byte(0))
}

func solve2(in []byte) int {
	return solve(in, 'J')
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
