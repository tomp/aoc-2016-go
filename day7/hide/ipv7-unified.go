package ipv7

import (
	'unicode/utf8'
)

type itemType int

const (
	itemError itemType := iota
	itemText
	itemSupernet
	itemHypernet
)

const (
	leftBracket := '['
	rightBacket := ']'
	netChars := 'abcdefghijklmnopqrstuvwxyz'
)

type item struct {
	typ itemType
	text string
}

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
    name  string    // used only for error reports.
    input string    // the string being scanned.
    start int       // start position of this item.
    pos   int       // current position in the input.
    width int       // width of last rune read from input.
    items chan item // channel of scanned items.
}

func lex(name, input string) (*lexer, chan item) {
    l := &lexer{
        name:  name,
        input: input,
        items: make(chan item),
    }
    go l.run()  // Concurrently run state machine.
    return l, l.items
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
    for state := lexText; state != nil; {
        state = state(l)
    }
    close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
    l.items <- item{t, l.input[l.start:l.pos]}
    l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() (rune int) {
    if l.pos >= len(l.input) {
        l.width = 0
        return eof
    }
    rune, l.width =
        utf8.DecodeRuneInString(l.input[l.pos:])
    l.pos += l.width
    return rune
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
    l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
    l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() int {
    rune := l.next()
    l.backup()
    return rune
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
    if strings.IndexRune(valid, l.next()) >= 0 {
        return true
    }
    l.backup()
    return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
    for strings.IndexRune(valid, l.next()) >= 0 {
    }
    l.backup()
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{})
  stateFn {
    l.items <- item{
        itemError,
        fmt.Sprintf(format, args...),
    }
    return nil
}


// lexText parses a non-IPv7 sequence
func lexText(l *lexer) stateFn {
	for {
		ch := l.input[l.pos]
        if isValidChar(ch) {
            if l.pos > l.start {
                l.emit(itemText)
            }
            return lexSupernet    //  Parse supernet
        } else if ch == rightBracket {
            if l.pos > l.start {
                l.emit(itemText)
            }
            return lexHypernet    // Parse hypernet
		}

        if l.next() == eof { break }
    }
    // Correctly reached EOF.
    if l.pos > l.start {
        l.emit(itemText)
    }
    l.emit(itemEOF)  // Useful to make EOF a token.
    return nil       // Stop the run loop.
}

// lexSupernet parses an IPv7 supernet sequence
func lexSupernet(l *lexer) stateFn {
	for {
		ch := l.input[l.pos]
        if ! isValidChar(ch) {
            if l.pos > l.start {
                l.emit(itemSupernet)
            }
			if ch == leftBracket {
				l.next()
				return lexHypernet    // Parse hypernet
			} else if ch == rightBracket {
				return errorf("Unexpected right bracket")
			} else {
				return lexText
			}
        }

        if l.next() == eof { break }
    }
    // Correctly reached EOF.
    if l.pos > l.start {
        l.emit(itemText)
    }
    l.emit(itemEOF)  // Useful to make EOF a token.
    return nil       // Stop the run loop.
}

// lexHypernet parses an IPv7 hypernet sequence
func lexHypernet(l *lexer) stateFn {
	for {
		if l.start == l.pos
		ch := l.input[l.pos]
        if ! isValidChar(ch) {
            if l.pos > l.start {
                l.emit(itemSupernet)
            }
			if ch == rightBracket {
				return lexHypernet    // Parse hypernet
			} else if ch == leftBracket {
				return errorf("Unexpected right bracket")
			} else {
				return lexText
			}
        }

        if l.next() == eof { break }
    }
    // Correctly reached EOF.
    if l.pos > l.start {
        l.emit(itemText)
    }
    l.emit(itemEOF)  // Useful to make EOF a token.
    return nil       // Stop the run loop.
}



