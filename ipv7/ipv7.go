package ipv7

import (
	"fmt"
	"github.com/tomp/aoc-2016-go/lexer"
	"strings"
)

const (
	ItemText lexer.ItemType = iota + 1
	ItemSupernet
	ItemHypernet
)

const (
	leftBracket   rune   = '['
	rightBracket  rune   = ']'
	leftBrackets  string = "["
	rightBrackets string = "]"
	netChars      string = "abcdefghijklmnopqrstuvwxyz"
)

type IPv7 struct {
	Addr      string
	Supernets []string
	Hypernets []string
}

func New(addr string) (ip IPv7, err error) {
	ip = IPv7{addr, nil, nil}
	ip.Supernets = make([]string, 0, 5)
	ip.Hypernets = make([]string, 0, 5)
	err = ip.parseAddr()
	return
}

func (ip *IPv7) parseAddr() (err error) {
	_, tokenChan := lexer.New("ipv7", ip.Addr, lexText)
	for item := range tokenChan {
		switch {
		case item.Typ == lexer.ItemEOF:
			break
		case item.Typ == lexer.ItemError:
			err = fmt.Errorf(item.Text)
			break
		case item.Typ == ItemSupernet:
			ip.Supernets = append(ip.Supernets, item.Text)
		case item.Typ == ItemHypernet:
			ip.Hypernets = append(ip.Hypernets, item.Text)
		}
	}
	return
}

func firstABBA(text string) (abba string) {
	for i := 0; i <= len(text)-4; i++ {
		if text[i] == text[i+3] &&
			text[i+1] == text[i+2] &&
			text[i] != text[i+1] {
			abba = text[i : i+4]
			break
		}
	}
	return
}

func allABA(text string) (abas []string) {
	for i := 0; i <= len(text)-3; i++ {
		if text[i] == text[i+2] &&
			text[i] != text[i+1] {
			abas = append(abas, text[i:i+3])
		}
	}
	return
}

func (ip *IPv7) IsTLS() bool {
	for _, hyper := range ip.Hypernets {
		if abba := firstABBA(hyper); abba != "" {
			return false
		}
	}
	for _, super := range ip.Supernets {
		if abba := firstABBA(super); abba != "" {
			return true
		}
	}
	return false
}

func (ip *IPv7) IsSSL() bool {
	all_aba := []string{}
	for _, super := range ip.Supernets {
		if abas := allABA(super); len(abas) > 0 {
			all_aba = append(all_aba, abas...)
		}
	}
	if len(all_aba) == 0 {
		return false
	}
	for _, aba := range all_aba {
		ab := strings.Split(aba[:2], "")
		bab := ab[1] + ab[0] + ab[1]
		for _, hyper := range ip.Hypernets {
			if strings.Index(hyper, bab) >= 0 {
				return true
			}
		}
	}
	return false
}

// lexText parses a non-IPv7 sequence
func lexText(l *lexer.State) lexer.StateFn {
	for {
		ch := l.Peek()
		switch {
		case strings.IndexRune(netChars, ch) >= 0:
			l.EmitIfToken(ItemText)
			return lexSupernet
		case ch == leftBracket:
			l.EmitIfToken(ItemText)
			l.Next()
			return lexHypernet
		case ch == rightBracket:
			l.EmitIfToken(ItemText)
			return l.Errorf("Unexpected right bracket")
		}

		if l.Next() == lexer.EOF {
			break
		}
	}
	// Correctly reached EOF.
	l.Emit(lexer.ItemEOF) // Useful to make EOF a token.
	return nil            // Stop the run loop.
}

// lexSupernet parses an IPv7 supernet sequence
func lexSupernet(l *lexer.State) lexer.StateFn {
	l.AcceptRun(netChars)
	l.EmitIfToken(ItemSupernet)

	ch := l.Peek()
	switch {
	case ch == leftBracket:
		return lexHypernet // Parse hypernet
	case ch == rightBracket:
		return l.Errorf("Unexpected right bracket")
	case ch == lexer.EOF:
		l.Emit(lexer.ItemEOF) // Useful to make EOF a token.
		return nil
	default:
		return lexText
	}
}

// lexHypernet parses an IPv7 hypernet sequence
func lexHypernet(l *lexer.State) lexer.StateFn {
	if !l.Accept(leftBrackets) {
		return lexText
	}
	l.AcceptRun(netChars)
	if !l.Accept(rightBrackets) {
		return l.Errorf("Missing right bracket in '%s'", l.Input())
	}
	l.Emit(ItemHypernet)

	ch := l.Peek()
	switch {
	case strings.IndexRune(netChars, ch) >= 0:
		return lexSupernet
	case ch == leftBracket:
		return lexHypernet
	case ch == lexer.EOF:
		l.Emit(lexer.ItemEOF) // Useful to make EOF a token.
		return nil
	default:
		return lexText
	}
}
