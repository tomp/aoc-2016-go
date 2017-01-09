package astar

import (
	"container/heap"
	"fmt"
)

type SearchState interface {
	String() string
	Hash() string
	AstarNextStates(visited map[string]bool) []SearchState
	Heuristic() int
	Done() bool
}

type SearchItem struct {
	priority int // lower priorities are considered first
	state    SearchState
	history  []SearchState
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

func Search(initState SearchState) (shortestPath []SearchState, err error) {
	queue := SearchQueue{}
	heap.Init(&queue)

	item := &SearchItem{priority: initState.Heuristic(),
		state:   initState,
		history: []SearchState{}}
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
		path := make([]SearchState, nsteps+1)
		for i, pastState := range item.history {
			path[i] = pastState
		}
		path[nsteps] = state

		if state.Done() {
			shortestPath = path
			break
		}

		count += 1
		for _, newState := range state.AstarNextStates(visited) {
			score := len(path) + newState.Heuristic()
			item = &SearchItem{priority: score,
				state:   newState,
				history: path}
			heap.Push(&queue, item)
		}
	}
	return
}

func PrintPath(path []SearchState) {
	for step, state := range path {
		fmt.Printf("\nStep %d\n", step)
		fmt.Println(state.String())
	}
	fmt.Println("** DONE **\n")
}
