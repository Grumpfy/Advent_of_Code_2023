package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var input []byte

type ModuleData struct {
	kind    byte
	childs  []string
	parents []string
}
type Modules map[string]*ModuleData

func parseInput(in []byte) Modules {
	modules := make(Modules)
	scanner := bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		line := scanner.Bytes()
		data := ModuleData{
			kind:    line[0],
			childs:  make([]string, 0, 8),
			parents: make([]string, 0, 8),
		}
		childs := line[7:]
		moduleName := string(line[1:3])
		if data.kind == 'b' {
			childs = line[15:]
			moduleName = "broadcaster"
		}
		for _, child := range bytes.Split(childs, []byte{',', ' '}) {
			data.childs = append(data.childs, string(child))
		}
		modules[moduleName] = &data
	}
	for k, v := range modules {
		for _, c := range v.childs {
			child := modules[c]
			if child != nil {
				child.parents = append(child.parents, k)
			}
		}
	}
	return modules
}

func solve1(in []byte) int {
	modules := parseInput(in)
	ffState := make(map[string]bool)
	cState := make(map[string]([]bool))

	for k, v := range modules {
		switch v.kind {
		case '%':
			ffState[k] = false
		case '&':
			cState[k] = make([]bool, len(v.parents))
		}
	}

	type Impulse struct {
		v    bool
		dest string
		from string
	}

	lowCount := 0
	highCount := 0

	for i := 0; i < 1000; i++ {
		queue := make([]Impulse, 0, 2048)
		lowCount++
		queue = append(queue, Impulse{false, "broadcaster", "button"})
		for len(queue) > 0 {
			imp := queue[0]
			queue = queue[1:]
			destModule := modules[imp.dest]
			if destModule == nil {
				continue
			}
			sendPulse := false
			pulseValue := false
			switch destModule.kind {
			case 'b':
				sendPulse = true
				pulseValue = imp.v
			case '%':
				currentState := ffState[imp.dest]
				newState := (currentState == imp.v)
				if newState != currentState {
					ffState[imp.dest] = newState
					sendPulse = true
					pulseValue = newState
				}
			case '&':
				state := cState[imp.dest]
				for i := range state {
					if destModule.parents[i] == imp.from {
						state[i] = imp.v
					}
				}
				pulse := true
				for _, p := range state {
					pulse = pulse && p
				}
				pulse = !pulse
				sendPulse = true
				pulseValue = pulse
			}
			if sendPulse {
				for _, c := range destModule.childs {
					if pulseValue {
						highCount++
					} else {
						lowCount++
					}
					queue = append(queue, Impulse{pulseValue, c, imp.dest})
				}
			}
		}
	}

	return lowCount * highCount
}

func solve2(in []byte) int {
	modules := parseInput(in)
	ffState := make(map[string]bool)
	cState := make(map[string]([]bool))

	for k, v := range modules {
		switch v.kind {
		case '%':
			ffState[k] = false
		case '&':
			cState[k] = make([]bool, len(v.parents))
		}
	}

	type Impulse struct {
		v    bool
		dest string
		from string
	}

	for i := 0; i < 10000; i++ {
		queue := make([]Impulse, 0, 2048)
		queue = append(queue, Impulse{false, "broadcaster", "button"})
		for len(queue) > 0 {
			imp := queue[0]
			queue = queue[1:]
			destModule := modules[imp.dest]
			if destModule == nil {
				continue
			}
			sendPulse := false
			pulseValue := false
			switch destModule.kind {
			case 'b':
				sendPulse = true
				pulseValue = imp.v
			case '%':
				currentState := ffState[imp.dest]
				newState := (currentState == imp.v)
				if newState != currentState {
					ffState[imp.dest] = newState
					sendPulse = true
					pulseValue = newState
				}
			case '&':
				state := cState[imp.dest]
				for i := range state {
					if destModule.parents[i] == imp.from {
						state[i] = imp.v
					}
				}
				pulse := true
				for _, p := range state {
					pulse = pulse && p
				}
				pulse = !pulse
				sendPulse = true
				pulseValue = pulse
			}
			if sendPulse {
				for _, c := range destModule.childs {
					if pulseValue && (c == "cn") {
						fmt.Println(imp.dest, i+1)
					}
					queue = append(queue, Impulse{pulseValue, c, imp.dest})
				}
			}
		}
	}
	return 3917*3943*3947*4001
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
