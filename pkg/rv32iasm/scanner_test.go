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
	auipc sp, 1`

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
