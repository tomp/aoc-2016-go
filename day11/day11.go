package main

import (
	"fmt"
	"github.com/tomp/aoc-2016-go/astar"
	"github.com/tomp/aoc-2016-go/rtg"
)

func main() {

	fmt.Println("## Example")

	state, err := rtg.InitialState(1, []int{2, 3}, []int{1, 1}, []string{"H", "L"})
	if err != nil {
		panic(err)
	}

	fmt.Println(state.String())

	path, err := astar.Search(&state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	astar.PrintPath(path)

	fmt.Println("\n## Part 1")

	state, err = rtg.InitialState(1, []int{1, 3, 3, 1, 1},
		[]int{2, 3, 3, 2, 1},
		[]string{"P", "Q", "R", "S", "T"})
	if err != nil {
		panic(err)
	}
	path, err = astar.Search(&state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	astar.PrintPath(path)

	fmt.Println("\n## Part 2")

	state, err = rtg.InitialState(1, []int{1, 3, 3, 1, 1, 1, 1},
		[]int{2, 3, 3, 2, 1, 1, 1},
		[]string{"P", "Q", "R", "S", "T", "D", "E"})
	if err != nil {
		panic(err)
	}
	path, err = astar.Search(&state)
	if err != nil {
		fmt.Println("ERROR: no solution found")
		return
	}
	astar.PrintPath(path)

}
