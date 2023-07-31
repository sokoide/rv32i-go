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
