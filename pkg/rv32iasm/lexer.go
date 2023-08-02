package rv32iasm

import (
	log "github.com/sirupsen/logrus"
)

const (
	// EOF end of file
	EOF = -1
	// UNKNOWN unknown token
	UNKNOWN = 0 * iota
)

var keywords = map[string]int{}

type token struct {
	tok int
	lit string
	pos position
}

type position struct {
	Line   int
	Column int
}

type lexer struct {
	s         *Scanner
	recentLit string
	recentPos position
	program   *Program
}

// Lex Called by goyacc
func (l *lexer) Lex(lval *assemblerSymType) int {
	tok, lit, pos, err := l.s.Scan()
	if err != nil {
		log.Errorf("%v", err)
	}
	if tok == EOF {
		return 0
	}
	lval.tok = token{tok: tok, lit: lit, pos: pos}
	l.recentLit = lit
	l.recentPos = pos
	return tok
}

// Error Called by goyacc
func (l *lexer) Error(e string) {
	log.Fatalf("Line %d, Column %d: %q %s",
		l.recentPos.Line, l.recentPos.Column, l.recentLit, e)
}
