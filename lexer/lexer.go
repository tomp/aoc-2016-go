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

func New(name, input string, initState StateFn) (*State, chan Item) {
	l := &State{
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

// emit passes an item back to the client.
func (l *State) Emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// emit passes an item back to the client.
func (l *State) EmitIfToken(t ItemType) {
	if l.pos > l.start {
		l.items <- Item{t, l.input[l.start:l.pos]}
		l.start = l.pos
	}
}

// next returns the next rune in the input.
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

// ignore skips over the pending input before this point.
func (l *State) Ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *State) Backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *State) Peek() (ch rune) {
	ch = l.Next()
	l.Backup()
	return ch
}

// accept consumes the next rune
// if it's from the valid set.
func (l *State) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *State) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

// acceptRunUntil consumes a run of runes up to the first occurrence of
// a run in the stope set.
func (l *State) AcceptRunUntil(stop string) {
	for strings.IndexRune(stop, l.Next()) < 0 {
	}
	l.Backup()
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *State) Errorf(format string, args ...interface{}) StateFn {
	l.items <- Item{
		ItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}
