// Package lexer provides the core code for a simple Rob-Pike-style lexer.
//
// To use this package, the client needs to define state functions using
// the given utility functions, and then initiate the state machine by
// calling New with the initial state and the string to be lexed.
//
// The client also needs to define item types for the tokens they wish to parse.
// Item types of value 0 or less are reserved for the lexer package.
//
// See lexer_test.go and ipv7/ipv7.go for examples.
//
package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ItemType int

const ItemError ItemType = 0
const ItemEOF ItemType = -1

const EOF rune = -1

type Item struct {
	Typ  ItemType
	Text string
}

type StateFn func(*State) StateFn

// State holds the state of the scanner.
type State struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan Item // channel of scanned items.
}

// New initializes and executes a state machine to lex the given input
// string.  A State object is returned, along with a read-only channel
// from which lexed token Items should be read.
func New(name, input string, initState StateFn) (State, chan Item) {
	l := State{
		name:  name,
		input: input,
		items: make(chan Item),
	}
	go l.run(initState) // Concurrently run state machine.
	return l, l.items
}

func (l *State) Input() string { return l.input }
func (l *State) Name() string  { return l.name }

// run lexes the input by executing state functions until
// the state is nil.
func (l *State) run(initState StateFn) {
	for state := initState; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// Emit passes an Item back to the client.  The item may be empty.
func (l *State) Emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// EmitIfToken passes an Item back to the client, if one has been found.
func (l *State) EmitIfToken(t ItemType) {
	if l.pos > l.start {
		l.items <- Item{t, l.input[l.start:l.pos]}
		l.start = l.pos
	}
}

// Next returns the next rune in the input.
func (l *State) Next() (ch rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	ch, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return ch
}

// Ignore skips over the pending input before this point.
func (l *State) Ignore() {
	l.start = l.pos
}

// Backup steps back one rune.
// Can only be used once per call of next.
func (l *State) Backup() {
	l.pos -= l.width
	l.width = 0
}

// Peek returns but does not consume
// the next rune in the input.
func (l *State) Peek() (ch rune) {
	ch = l.Next()
	l.Backup()
	return ch
}

// Accept consumes the next rune
// if it's from the valid set.
func (l *State) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *State) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

// AcceptRunUntil consumes a run of runes up to the first occurrence of
// a run in the stop set.
func (l *State) AcceptRunUntil(stop string) {
	ch := l.Peek()
	for ; ch != EOF && strings.IndexRune(stop, ch) < 0; ch = l.Next() {
	}
	l.Backup()
}

// Errorf returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *State) Errorf(format string, args ...interface{}) StateFn {
	l.items <- Item{
		ItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}
