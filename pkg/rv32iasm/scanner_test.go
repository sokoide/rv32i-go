package rv32iasm

import (
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {
	src := `boot:
# This is a comment line
	lui a0, 4 # This is a comment
	auipc sp, 1
	beq ra, a0, 123
	bne ra, a0, 123
	blt ra, a0, 123
	bge ra, a0, 123
	bltu ra, a0, 123
	bgeu ra, a0, 123
	lb ra, -100(a0)
	lh ra, -100(a0)
	lw ra, -100(a0)
	lbu ra, -100(a0)
	lhu ra, -100(a0)
	sb ra, -100(a0)
	sh ra, -100(a0)
	sw ra, -100(a0)
	addi ra, a0, -123
	li ra, 0
	li s0, 0
	slti ra, a0, -123
	sltiu ra, a0, 123
	seqz ra, a0
	xori ra, a0, -123
	ori ra, a0, -123
	andi ra, a0, -123
	slli ra, a0, 123
	srli ra, a0, 123
	srai ra, a0, 123
# Another comment`

	reader := strings.NewReader(src)
	scanner := NewScanner(reader)
	program, err := scanner.Parse()
	if err != nil {
		t.Error(err)
	}

	wants := []statement{
		{"label", 0, 0, 0, "boot"},
		{"comment", 0, 0, 0, ""},
		{"lui", 10, 4, 0, ""},
		{"auipc", 2, 1, 0, ""},
		{"beq", 1, 10, 123, ""},
		{"bne", 1, 10, 123, ""},
		{"blt", 1, 10, 123, ""},
		{"bge", 1, 10, 123, ""},
		{"bltu", 1, 10, 123, ""},
		{"bgeu", 1, 10, 123, ""},
		{"lb", 1, -100, 10, ""},
		{"lh", 1, -100, 10, ""},
		{"lw", 1, -100, 10, ""},
		{"lbu", 1, -100, 10, ""},
		{"lhu", 1, -100, 10, ""},
		{"sb", 1, -100, 10, ""},
		{"sh", 1, -100, 10, ""},
		{"sw", 1, -100, 10, ""},
		{"addi", 1, 10, -123, ""},
		{"li", 1, 0, 0, ""},
		{"li", 8, 0, 0, ""},
		{"slti", 1, 10, -123, ""},
		{"sltiu", 1, 10, 123, ""},
		{"seqz", 1, 10, 0, ""},
		{"xori", 1, 10, -123, ""},
		{"ori", 1, 10, -123, ""},
		{"andi", 1, 10, -123, ""},
		{"slli", 1, 10, 123, ""},
		{"srli", 1, 10, 123, ""},
		{"srai", 1, 10, 123, ""},
		{"comment", 0, 0, 0, ""},
	}

	if len(program.statements) != len(wants) {
		t.Error("Wrong length")
	}

	for idx, got := range program.statements {
		want := wants[idx]
		if *got != want {
			t.Errorf("[%d] got:%+v, want:%+v", idx, got, want)
		}
	}
}
