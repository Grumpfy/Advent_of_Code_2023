package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"time"
)

//go:embed input.txt
var input []byte

var partProp = [256]uint8{'x': 0, 'm': 1, 'a': 2, 's': 3}

type Part [4]int

type Rule struct {
	output    string
	op        byte
	partId    int
	threshold int
}

type Rules []Rule

type Workflows map[string]Rules

func parseInput(in []byte) (Workflows, []Part) {
	scanner := bufio.NewScanner(bytes.NewReader(in))
	workflows := make(Workflows)
	parts := make([]Part, 0, 1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if line[0] != '{' {
			// ex: cms{s>621:A,x<502:dz,m<3295:R,nzd}
			bracketPos := bytes.IndexByte(line, '{')
			identifier := line[0:bracketPos]
			rulesInput := line[bracketPos+1 : len(line)-1]
			rules := make(Rules, 0, 8)
			for _, ruleInput := range bytes.Split(rulesInput, []byte{','}) {
				colPos := bytes.IndexByte(ruleInput, ':')
				if colPos == -1 {
					rules = append(rules, Rule{string(ruleInput), 0, 0, 0})
				} else {
					r, _ := strconv.Atoi(string(ruleInput[2:colPos]))
					rules = append(rules, Rule{string(ruleInput[colPos+1:]), ruleInput[1], int(partProp[ruleInput[0]]), r})
				}
			}
			workflows[string(identifier)] = rules
		} else {
			var part Part
			for i, prop := range bytes.Split(line[1:len(line)-1], []byte{','}) {
				part[i], _ = strconv.Atoi(string(prop[2:]))
			}
			parts = append(parts, part)
		}
	}
	return workflows, parts
}

func solve1(in []byte) int {
	workflows, parts := parseInput(in)
	res := 0
	for _, part := range parts {
		workflow := "in"
		for {
			if workflow == "R" {
				break
			}
			if workflow == "A" {
				res += part[0] + part[1] + part[2] + part[3]
				break
			}
			for _, rule := range workflows[workflow] {
				ok := false
				switch rule.op {
				case '<':
					if part[rule.partId] < rule.threshold {
						ok = true
					}
				case '>':
					if part[rule.partId] > rule.threshold {
						ok = true
					}
				default:
					ok = true
				}
				if ok {
					workflow = rule.output
					break
				}
			}
		}
	}

	return res
}

type Range [2]int
type Ranges [4]Range

func solveRanges(workflows *Workflows, ranges Ranges, input string) int {
	if input == "R" {
		return 0
	}
	if input == "A" {
		res := 1
		for _, r := range ranges {
			res *= r[1] - r[0] + 1
		}
		return res
	}

	res := 0
	for _, rule := range (*workflows)[input] {
		switch rule.op {
		case '<':
			if ranges[rule.partId][0] < rule.threshold {
				passRange := ranges
				passRange[rule.partId][1] = min(rule.threshold-1, ranges[rule.partId][1])
				res += solveRanges(workflows, passRange, rule.output)
				ranges[rule.partId][0] = rule.threshold
			}
		case '>':
			if ranges[rule.partId][1] > rule.threshold {
				passRange := ranges
				passRange[rule.partId][0] = max(rule.threshold+1, ranges[rule.partId][0])
				res += solveRanges(workflows, passRange, rule.output)
				ranges[rule.partId][1] = rule.threshold
			}
		default:
			res += solveRanges(workflows, ranges, rule.output)
		}
		if ranges[rule.partId][0] > ranges[rule.partId][1] {
			break
		}
	}

	return res
}

func solve2(in []byte) int {
	workflows, _ := parseInput(in)
	ranges := Ranges{Range{1, 4000}, Range{1, 4000}, Range{1, 4000}, Range{1, 4000}}
	return solveRanges(&workflows, ranges, "in")
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
