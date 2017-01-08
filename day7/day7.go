package main

import (
	"bufio"
	"fmt"
	"github.com/tomp/aoc-2016-go/ipv7"
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

	tlsCount := 0
	for _, line := range lines {
		ip, err := ipv7.New(line)
		if err == nil && ip.IsTLS() {
			tlsCount += 1
		}
	}
	fmt.Printf("%d TLS addresses found\n", tlsCount)

	sslCount := 0
	for _, line := range lines {
		ip, err := ipv7.New(line)
		if err == nil && ip.IsSSL() {
			sslCount += 1
		}
	}
	fmt.Printf("%d SSL addresses found\n", sslCount)
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
