package main

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

type lexer struct {
	s         *scanner
	recentLit string
	recentPos position
	program   *program
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

func parse(s *scanner) *program {
	l := lexer{s: s}
	l.program = &program{
		statements: make([]*statement, 0),
	}
	if assemblerParse(&l) != 0 {
		panic("Parse error")
	}
	return l.program
}

func main() {
	ev := NewEvaluator()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	// log.SetLevel(log.InfoLevel)

	log.Info("asm started")
	src := `boot:
# This is a comment line
	li ra, -300
	li s0, 0 # This is a comment
	lui a0, 4
label2:
	li a1, 1000000000`

	log.Tracef("src: %s", src)

	s := bufio.NewScanner(strings.NewReader(src))
	scanner := new(scanner)
	source := []string{}
	for s.Scan() {
		source = append(source, s.Text())
	}
	scanner.Init(strings.Join(source, "\n") + "\n")

	var program *program = parse(scanner)
	log.Debugf("* program=%+v", program)

	log.Info("* start evaluation")
	err := ev.evaluate_program(program)
	if err != nil {
		panic(nil)
	}

	emu := rv32i.NewEmulator()
	emu.Reset()
}
