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
	slti ra, a0, -123
	sltiu ra, a0, 123
	xori ra, a0, -123
	ori ra, a0, -123
	andi ra, a0, -123
	slli ra, a0, 123
	srli ra, a0, 123
	srai ra, a0, 123
	add ra, a0, a1
	sub ra, a0, a1
	sll ra, a0, a1
	slt ra, a0, a1
	sltu ra, a0, a1
	xor ra, a0, a1
	srl ra, a0, a1
	sra ra, a0, a1
	or ra, a0, a1
	and ra, a0, a1
	call a0, hoge
	call hoge
	li ra, 0
	li s0, 0
	seqz ra, a0
	ret
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
		{"slti", 1, 10, -123, ""},
		{"sltiu", 1, 10, 123, ""},
		{"xori", 1, 10, -123, ""},
		{"ori", 1, 10, -123, ""},
		{"andi", 1, 10, -123, ""},
		{"slli", 1, 10, 123, ""},
		{"srli", 1, 10, 123, ""},
		{"srai", 1, 10, 123, ""},
		{"add", 1, 10, 11, ""},
		{"sub", 1, 10, 11, ""},
		{"sll", 1, 10, 11, ""},
		{"slt", 1, 10, 11, ""},
		{"sltu", 1, 10, 11, ""},
		{"xor", 1, 10, 11, ""},
		{"srl", 1, 10, 11, ""},
		{"sra", 1, 10, 11, ""},
		{"or", 1, 10, 11, ""},
		{"and", 1, 10, 11, ""},
		{"call", 10, 0, 0, "hoge"},
		{"call", 1, 0, 0, "hoge"},
		{"li", 1, 0, 0, ""},
		{"li", 8, 0, 0, ""},
		{"sltiu", 1, 10, 1, ""},
		{"jalr", 0, 0, 1, ""},
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
