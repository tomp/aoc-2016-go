package main

import (
	"bufio"
	"fmt"
	"github.com/tomp/aoc-2016-go/asmbunny"
	"os"
)

const (
	INPUTFILE string = "input.txt"
)

func main() {
	// Part1
	lines, err := read_lines(INPUTFILE)
	if err != nil {
		panic(err)
	}

	prog, err := asmbunny.Compile(lines)
	if err != nil {
		panic(err)
	}

	fmt.Println("## Part 1")
	init := asmbunny.Registers{}
	reg, err := prog.Execute(init)

	fmt.Printf("a: %d \n", reg.Get("a"))
	fmt.Printf("b: %d \n", reg.Get("b"))
	fmt.Printf("c: %d \n", reg.Get("c"))
	fmt.Printf("d: %d \n", reg.Get("d"))

	fmt.Println("\n## Part 2")
	init = asmbunny.Registers{}
	init.Set("c", 1)
	reg, err = prog.Execute(init)

	fmt.Printf("a: %d \n", reg.Get("a"))
	fmt.Printf("b: %d \n", reg.Get("b"))
	fmt.Printf("c: %d \n", reg.Get("c"))
	fmt.Printf("d: %d \n", reg.Get("d"))
}

// read_lines returns the contents of the given file as a slice
// of lines.
func read_lines(filename string) (lines []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
