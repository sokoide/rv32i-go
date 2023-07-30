package main

import (
	"errors"
)

const (
	// EOF end of file
	EOF = -1
	// UNKNOWN unknown token
	UNKNOWN = 0 * iota
)

var keywords = map[string]int{}

var regs = map[string]int{
	"zero": 0,
	"ra":   1,
	"sp":   2,
	"gp":   3,
	"tp":   4,
	"t0":   5,
	"t1":   6,
	"t2":   7,
	"s0":   8,
	"fp":   8,
	"s1":   9,
	"a0":   10,
	"a1":   11,
	"a2":   12,
	"a3":   13,
	"a4":   14,
	"a5":   15,
	"a6":   16,
	"a7":   17,
	"s2":   18,
	"s3":   19,
	"s4":   20,
	"s5":   21,
	"s6":   22,
	"s7":   23,
	"s8":   24,
	"s9":   25,
	"s10":  26,
	"s11":  27,
	"t3":   28,
	"t4":   29,
	"t5":   30,
	"t6":   31,
	"x0":   0,
	"x1":   1,
	"x2":   2,
	"x3":   3,
	"x4":   4,
	"x5":   5,
	"x6":   6,
	"x7":   7,
	"x8":   8,
	"x9":   9,
	"x10":  10,
	"x11":  11,
	"x12":  12,
	"x13":  13,
	"x14":  14,
	"x15":  15,
	"x16":  16,
	"x17":  17,
	"x18":  18,
	"x19":  19,
	"x20":  20,
	"x21":  21,
	"x22":  22,
	"x23":  23,
	"x24":  24,
	"x25":  25,
	"x26":  26,
	"x27":  27,
	"x28":  28,
	"x29":  29,
	"x30":  30,
	"x31":  31,
}

type token struct {
	tok int
	lit string
	pos position
}

type position struct {
	Line   int
	Column int
}

type scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *scanner) Init(src string) {
	s.src = []rune(src)
}

func (s *scanner) Scan() (tok int, lit string, pos position, err error) {
	s.skipWhiteSpace()
	s.skipComment()
	pos = s.position()
	switch ch := s.peek(); {
	case isDigit(ch):
		tok, lit = NUMBER, s.scanNumber()
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = s.tokFromLit(lit)
	case ch == '-':
		s.next()
		if !isDigit(s.peek()) {
			err = errors.New("syntax error. it should have number(s) after '-")
		}
		tok, lit = NUMBER, "-"+s.scanNumber()
	case ch == '\n':
		s.next()
		tok, lit = LF, ""
	case ch == ':':
		s.next()
		tok, lit = COLON, ":"
	case ch == ',':
		s.next()
		tok, lit = COMMA, ","
	default:
		switch ch {
		case -1:
			tok = EOF
		case '(', ')', ';', '+', '-', '*', '/', '%', '=':
			tok = int(ch)
			lit = string(ch)
		}
		s.next()
	}
	return
}

// ========================================

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func (s *scanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *scanner) position() position {
	return position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *scanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *scanner) skipComment() {
	if s.peek() == '#' {
		s.next()
		for s.peek() != '\n' {
			s.next()
		}
	}
}

func (s *scanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}

func (s *scanner) tokFromLit(lit string) int {
	if _, ok := regs[lit]; ok {
		return REGISTER
	}

	switch lit {
	case "li":
		return LI
	case "lui":
		return LUI
	default:
		return IDENT
	}
}

func (s *scanner) scanNumber() string {
	var ret []rune
	for isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}
