package rv32iasm

import (
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {
	src := `boot:
# This is a comment line
	li ra, 0
	li s0, 0 # This is a comment
	lui a0, 4
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
		{"li", 1, 0, 0, ""},
		{"li", 8, 0, 0, ""},
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
