package lexer

import (
	"fmt"
	"testing"
)

const whitespace string = " \t"

const ItemWord ItemType = ItemError + 1

// Parse whitespace-delimited tokens
func lexWord(l *State) StateFn {
	l.AcceptRun(whitespace)
	l.Ignore()
	l.AcceptRunUntil(whitespace)
	if l.Peek() == EOF {
		l.EmitIfToken(ItemWord)
		l.Emit(ItemEOF)
		return nil
	}
	l.Emit(ItemWord)
	return lexWord
}

func TestLexer(t *testing.T) {
	cases := [...]struct {
		input string
		words []string
	}{
		{"a b c", []string{"a", "b", "c"}},
		{"  aaa\tbb \t c  ", []string{"aaa", "bb", "c"}},
	}

	for n, item := range cases {
		name := fmt.Sprintf("case %d", n)
		_, tokenChan := New(name, item.input, lexWord)

		results := []string{}
		for token := range tokenChan {
			if token.Typ == ItemEOF {
				break
			}
			results = append(results, token.Text)
		}
		if len(results) != len(item.words) {
			t.Errorf("%d words found in '%s'  (expected %d)",
				len(results), item.input, len(item.words))
		} else {
			for i, expected := range item.words {
				if results[i] != expected {
					t.Errorf("Word %d was '%s'  (expected '%s')",
						i, results[i], expected)
				}
			}
		}
	}
}
