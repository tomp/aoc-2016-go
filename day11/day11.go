package main

import (
	"container/heap"
	"fmt"
	"github.com/tomp/aoc-2016-go/rtg"
)

type SearchItem struct {
	priority int // lower priorities are considered first
	state    *rtg.State
	history  []*rtg.State
	index    int
}

type SearchQueue []*SearchItem

func (q SearchQueue) Len() int { return len(q) }

func (q SearchQueue) Less(i, j int) bool {
	return q[i].priority < q[j].priority
}

func (q SearchQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *SearchQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*SearchItem)
	item.index = n
	*q = append(*q, item)
}

func (q *SearchQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

func main() {

	fmt.Println("## Example")

	state, err := rtg.InitialState(1, []int{2, 3}, []int{1, 1}, []string{"H", "L"})
	if err != nil {
		panic(err)
	}

	fmt.Println(state.String())

	path, err := search(state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	printPath(path)

	fmt.Println("\n## Part 1")

	state, err = rtg.InitialState(1, []int{1, 3, 3, 1, 1},
		[]int{2, 3, 3, 2, 1},
		[]string{"P", "Q", "R", "S", "T"})
	if err != nil {
		panic(err)
	}
	path, err = search(state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	printPath(path)

	fmt.Println("\n## Part 2")

	state, err = rtg.InitialState(1, []int{1, 3, 3, 1, 1, 1, 1},
		[]int{2, 3, 3, 2, 1, 1, 1},
		[]string{"P", "Q", "R", "S", "T", "D", "E"})
	if err != nil {
		panic(err)
	}
	path, err = search(state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	printPath(path)

}

func printPath(path []*rtg.State) {
	for step, state := range path {
		fmt.Printf("\nStep %d\n", step)
		fmt.Println(state.String())
	}
	fmt.Println("** DONE **\n")
}

func search(initState rtg.State) (shortestPath []*rtg.State, err error) {
	queue := SearchQueue{}
	heap.Init(&queue)

	item := &SearchItem{priority: initState.Heuristic(),
		state:   &initState,
		history: []*rtg.State{}}
	heap.Push(&queue, item)

	visited := make(map[string]bool, 10)
	visited[initState.Hash()] = true

	count := 0 // number of states expanded
	for {
		item := heap.Pop(&queue).(*SearchItem)
		if err != nil {
			break
		}
		state := item.state

		// Add current state to history before generating next steps
		nsteps := len(item.history)
		path := make([]*rtg.State, nsteps+1)
		for i, pastState := range item.history {
			path[i] = pastState
		}
		path[nsteps] = state

		if state.Done() {
			shortestPath = path
			break
		}

		count += 1
		for _, newState := range state.NextStates(visited) {
			score := len(path) + newState.Heuristic()
			item = &SearchItem{priority: score,
				state:   newState,
				history: path}
			heap.Push(&queue, item)
		}
	}
	fmt.Printf("\n## %d states considered\n", count)
	return
}
