package main

import (
	"crypto/md5"
	"fmt"
	"container/heap"
	pq "github.com/tomp/aoc-2016-go/pqueue"
	"strings"
)

const (
	INPUT = "pgflpeqp"
	XSIZE = 4
	YSIZE = 4
)

type Direction string

const (
	UP    = "U"
	RIGHT = "R"
	DOWN  = "D"
	LEFT  = "L"
)

const OPEN_CODES = "bcdef"

// openDoors returns a string that reports the locations (Directions)
// of the open doors in the current position, given the specified
// passcode and sequence of steps taken so far.
func openDoors(passcode string, path History) (doors string) {
	x, y := position(path)
	key := passcode + path.steps
	hash := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	if strings.Contains(OPEN_CODES, string(hash[0])) && y > 0 {
		doors += UP
	}
	if strings.Contains(OPEN_CODES, string(hash[1])) && y < YSIZE-1 {
		doors += DOWN
	}
	if strings.Contains(OPEN_CODES, string(hash[2])) && x > 0 {
		doors += LEFT
	}
	if strings.Contains(OPEN_CODES, string(hash[3])) && x < XSIZE-1 {
		doors += RIGHT
	}
	return
}

// position calculates the final x, y coordinates for a unit
// starting at (0, 0) and following the given sequence of steps.
func position(path History) (x, y int) {
	for _, d := range strings.Split(path.steps, "") {
		switch d {
		case UP:
			y -= 1
		case RIGHT:
			x += 1
		case DOWN:
			y += 1
		case LEFT:
			x -= 1
		}
	}
	return
}

type History struct {
	steps string
}

func (h *History) String() string {
	return h.steps
}

func (h History) addStep(door string) History {
	return History{h.steps + door}
}

func (h History) Heuristic() int {
	x, y := position(h)
	return (XSIZE - x - 1) + (YSIZE - y - 1)
}

func search(passcode string, longest bool) (path string, nstates int) {
    queue := pq.Queue{}
	heap.Init(&queue)
	initState := History{}
	initScore := 0 + initState.Heuristic()
	heap.Push(&queue, pq.NewItem(initScore, &initState))
	longestPath := ""
	for len(queue) > 0 {
		nstates += 1
		item := heap.Pop(&queue).(*(pq.Item))
		state := *(item.Value.(*History))
		score := state.Heuristic()
		if score == 0 {
			if ! longest {
				return state.steps, nstates
			}
			if len(state.steps) > len(longestPath) {
				longestPath = state.steps
			}
			continue
		}
		for _, door := range openDoors(passcode, state) {
			newState := state.addStep(string(door))
			newScore := len(state.steps) + newState.Heuristic()
			heap.Push(&queue, pq.NewItem(newScore, &newState))
		}
	}
	// if we run out of states to explore without having found a
	// solution, then longestPath will be the empty string.
	path = longestPath
	return
}

func main() {

	fmt.Println("## Example")

	examples := [...]struct {
		passcode,
		solution string
		longest int
	} {
		{"ihgpwlah", "DDRRRD", 370},
		{"kglvqrro", "DDUDRLRRUDRD", 492},
		{"ulqzkmiv", "DRURDRUDDLLDLUURRDULRLDUUDDDRR", 830},
	}

	fmt.Println("Shortest path...")
	for _, item := range examples {
		path, nstates := search(item.passcode, false)
		fmt.Printf("Passcode: '%s'  Solution: %d steps  (%d states considered)\n",
			item.passcode, len(path), nstates)
		if path != item.solution {
			fmt.Printf("ERROR: incorrect solution found for passcode '%s'\n",
				item.passcode)
			fmt.Printf("Result: '%s', expected '%s'\n", path, item.solution)
			return
		}
	}

	fmt.Println("\nLongest path...")
	for _, item := range examples {
		path, nstates := search(item.passcode, true)
		fmt.Printf("Passcode: '%s'  Solution: %d steps  (%d states considered)\n",
			item.passcode, len(path), nstates)
		if len(path) != item.longest {
			fmt.Printf("ERROR: longest solution not found for passcode '%s'\n",
				item.passcode)
			fmt.Printf("Result: %d, expected %d\n", len(path), item.longest)
			return
		}
	}

	fmt.Println("\n## Part 1")
	fmt.Println("Shortest path...")

	passcode := "pgflpeqp"
	expected := "RDRLDRDURD"

	path, nstates := search(passcode, false)
	fmt.Printf("Passcode: '%s'  Solution: %d steps  (%d states considered)\n",
		passcode, len(path), nstates)
	fmt.Printf("Result: '%s'\n", path)
	if path != expected {
		fmt.Printf("D'oh!  The result should have been '%s'\n", expected)
		return
	}

	fmt.Println("\n## Part 2")
	fmt.Println("Longest path...")

	expectedLongest := 596

	path, nstates = search(passcode, true)
	fmt.Printf("Passcode: '%s'  Solution: %d steps  (%d states considered)\n",
		passcode, len(path), nstates)
	if len(path) != expectedLongest {
		fmt.Printf("D'oh!  The result should have been %d steps\n",
		    expectedLongest)
		return
	}
}
