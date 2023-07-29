package main

import "golang.org/x/exp/slices"

const (
	// EOF end of file
	EOF = -1
	// UNKNOWN unknown token
	UNKNOWN = 0 * iota
)

var keywords = map[string]int{}

// TODO: make it a map and map "sp" -> 2
var registers = []string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2", "s0", "fp", "s1", "a0",
	"a1", "a2", "a3", "a4", "a5", "a6", "a7", "s2", "s3", "s4", "s5", "s6",
	"s7", "s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
	"x0", "x1", "x2", "x3", "x4", "x5", // TODO:
}

var regs = map[string]int{
	"zero": 0,
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

func (s *scanner) Scan() (tok int, lit string, pos position) {
	s.skipWhiteSpace()
	s.skipComment()
	pos = s.position()
	switch ch := s.peek(); {
	case isDigit(ch):
		tok, lit = NUMBER, s.scanNumber()
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = s.tokFromLit(lit)
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
	if slices.Contains(registers, lit) {
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
