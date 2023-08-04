package rv32iasm

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func NewScanner(reader io.Reader) *Scanner {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	scanner := new(Scanner)
	scanner.init(buf.String())
	return scanner
}

func (s *Scanner) init(src string) {
	if !strings.HasSuffix(src, "\n") {
		src += "\n"
	}
	s.src = []rune(src)
}

func (s *Scanner) Parse() (*Program, error) {
	l := lexer{s: s}
	l.program = &Program{
		statements: make([]*statement, 0),
	}
	if assemblerParse(&l) != 0 {
		return nil, errors.New("Parse error")
	}
	return l.program, nil
}

func (s *Scanner) Scan() (tok int, lit string, pos position, err error) {
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
	case ch == '(':
		s.next()
		tok, lit = LP, "("
	case ch == ')':
		s.next()
		tok, lit = RP, ")"
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

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func (s *Scanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *Scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *Scanner) position() position {
	return position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *Scanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *Scanner) skipComment() {
	if s.peek() == '#' {
		s.next()
		for s.peek() != '\n' {
			s.next()
		}
	}
}

func (s *Scanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}

func (s *Scanner) tokFromLit(lit string) int {
	if _, ok := rv32i.Regs[lit]; ok {
		return REGISTER
	}

	switch lit {
	case "lui":
		return LUI
	case "auipc":
		return AUIPC
	case "jal":
		return JAL
	case "jalr":
		return JALR
	case "beq":
		return BEQ
	case "bne":
		return BNE
	case "blt":
		return BLT
	case "bge":
		return BGE
	case "bltu":
		return BLTU
	case "bgeu":
		return BGEU
	case "lb":
		return LB
	case "lh":
		return LH
	case "lw":
		return LW
	case "lbu":
		return LBU
	case "lhu":
		return LHU
	case "sb":
		return SB
	case "sh":
		return SH
	case "sw":
		return SW
	case "addi":
		return ADDI
	case "li":
		return LI
	case "slti":
		return SLTI
	case "sltiu":
		return SLTIU
	case "seqz":
		return SEQZ
	case "xori":
		return XORI
	case "ori":
		return ORI
	case "andi":
		return ANDI
	case "slli":
		return SLLI
	case "srli":
		return SRLI
	case "srai":
		return SRAI
	case "add":
		return ADD
	case "sub":
		return SUB
	case "sll":
		return SLL
	case "slt":
		return SLT
	case "sltu":
		return SLTU
	case "xor":
		return XOR
	case "srl":
		return SRL
	case "sra":
		return SRA
	case "or":
		return OR
	case "and":
		return AND
	case "ret":
		return RET
	default:
		return IDENT
	}
}

func (s *Scanner) scanNumber() string {
	var ret []rune
	for isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}
