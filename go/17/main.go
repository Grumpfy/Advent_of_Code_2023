package main

import (
	"bytes"
	"container/heap"
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

type Key struct {
	pos   int
	dir   uint8
	steps int
}

type Node struct {
	key       Key
	heatLoss  int
	heapIndex int
	prev      *Node
}

type MinHeap []*Node
type Nodes map[Key]*Node

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
	return h[i].heatLoss < h[j].heatLoss
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *MinHeap) Push(x any) {
	n := len(*h)
	item := x.(*Node)
	item.heapIndex = n
	*h = append(*h, item)
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.heapIndex = -1
	*h = old[0 : n-1]
	return item
}

func (n *Nodes) getOrCreate(key Key) (*Node, bool) {
	node, ok := (*n)[key]
	if ok {
		return node, false
	} else {
		newNode := &Node{}
		newNode.key = key
		(*n)[key] = newNode
		return newNode, true
	}
}

func childKey(key Key, dir uint8, width int, end int) (Key, error) {
	childIdx := -1
	switch dir {
	case Right:
		if (key.pos % width) < (width - 2) {
			childIdx = key.pos + 1
		}
	case Down:
		next := key.pos + width
		if next < end {
			childIdx = next
		}
	case Left:
		if (key.pos % width) > 0 {
			childIdx = key.pos - 1
		}
	case Up:
		next := key.pos - width
		if next >= 0 {
			childIdx = next
		}
	}

	if childIdx == -1 {
		return Key{}, errors.New("no child")
	}

	if dir == key.dir {
		return Key{childIdx, dir, key.steps + 1}, nil
	} else {
		return Key{childIdx, dir, 1}, nil
	}
}

func rotateL(dir uint8) uint8 {
	return (dir - 1) % 4
}

func rotateR(dir uint8) uint8 {
	return (dir + 1) % 4
}

func solve1(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	nodes := make(Nodes)
	priorityQ := make(MinHeap, 0, 1024)
	outNodePos := len(in) - 2

	startNode, _ := nodes.getOrCreate(Key{})
	priorityQ = append(priorityQ, startNode)

	for priorityQ.Len() != 0 {
		node := heap.Pop(&priorityQ).(*Node)
		if node.key.pos == outNodePos {
			return node.heatLoss
		}

		dirs := [3]uint8{rotateL(node.key.dir), node.key.dir, rotateR(node.key.dir)}
		for _, dir := range dirs {
			if key, err := childKey(node.key, dir, width, len(in)); err == nil && key.steps < 4 {
				childNode, created := nodes.getOrCreate(key)
				heatLoss := node.heatLoss + int(in[childNode.key.pos]-'0')
				if created {
					childNode.heatLoss = heatLoss
					childNode.prev = node
					heap.Push(&priorityQ, childNode)
				} else if childNode.heapIndex != -1 && heatLoss < childNode.heatLoss {
					childNode.heatLoss = heatLoss
					childNode.prev = node
					heap.Fix(&priorityQ, childNode.heapIndex)
				}
			}
		}
	}

	return -1
}

func solve2(in []byte) int {
	width := bytes.IndexByte(in, '\n') + 1
	nodes := make(Nodes)
	priorityQ := make(MinHeap, 0, 1024)
	outNodePos := len(in) - 2

	startNode1, _ := nodes.getOrCreate(Key{dir:0})
	startNode2, _ := nodes.getOrCreate(Key{dir:1})
	startNode1.heapIndex, startNode2.heapIndex = 0, 1
	priorityQ = append(priorityQ, startNode1, startNode2)
	heap.Init(&priorityQ)

	for priorityQ.Len() != 0 {
		node := heap.Pop(&priorityQ).(*Node)
		if node.key.pos == outNodePos && node.key.steps >= 4 {
			return node.heatLoss
		}

		dirs := [3]uint8{rotateL(node.key.dir), node.key.dir, rotateR(node.key.dir)}
		for i, dir := range dirs {
			if i != 1 && node.key.steps < 4 {
				continue
			}
			if i == 1 && node.key.steps >= 10 {
				continue
			}

			if key, err := childKey(node.key, dir, width, len(in)); err == nil {
				childNode, created := nodes.getOrCreate(key)
				heatLoss := node.heatLoss + int(in[childNode.key.pos]-'0')
				if created {
					childNode.heatLoss = heatLoss
					childNode.prev = node
					heap.Push(&priorityQ, childNode)
				} else if childNode.heapIndex != -1 && heatLoss < childNode.heatLoss {
					childNode.heatLoss = heatLoss
					childNode.prev = node
					heap.Fix(&priorityQ, childNode.heapIndex)
				}
			}
		}
	}

	return -1
}

func main() {
	start := time.Now()
	res := solve2(input)
	elapsed := time.Since(start)
	fmt.Println(elapsed, res)
}
