package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed input.txt
var input []byte

func solve1(in []byte) int {
	sum := 0
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Text()
		first := strings.IndexAny(line, "0123456789")
		last := strings.LastIndexAny(line, "0123456789")
		sum += 10*int(line[first]-'0') + int(line[last]-'0')
	}
	return sum
}

func solve2(in []byte) int {
	sum := 0
	digitNames := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Text()
		firstDigitIndex := strings.IndexAny(line, "0123456789")
		lastDigitIndex := strings.LastIndexAny(line, "0123456789")
		firstDigit := int(line[firstDigitIndex] - '0')
		lastDigit := int(line[lastDigitIndex] - '0')
		for n, digitName := range digitNames {
			if index := strings.Index(line, digitName); index != -1 {
				if index < firstDigitIndex {
					firstDigitIndex = index
					firstDigit = n + 1
				}
			}
			if index := strings.LastIndex(line, digitName); index != -1 {
				if index > lastDigitIndex {
					lastDigitIndex = index
					lastDigit = n + 1
				}
			}
		}
		sum += 10*firstDigit + lastDigit
	}
	return sum
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
