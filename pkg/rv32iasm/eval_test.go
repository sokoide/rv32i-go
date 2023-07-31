package rv32iasm

import (
	"strings"
	"testing"
)

func Test_EvaluateProgram(t *testing.T) {
	src := `boot:
# This is a comment line
	li ra, 0
	li s0, 0 # This is a comment
	lui a0, 4
	auipc sp, 1`

	reader := strings.NewReader(src)
	scanner := NewScanner(reader)
	program := scanner.Parse()

	ev := NewEvaluator()
	err := ev.EvaluateProgram(program)
	if err != nil {
		t.Error("Failed to evaluate")
	}

	wants := []uint32{
		0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410133,
	}
	if len(ev.Code) == len(wants) {
		t.Errorf("Unexpected length. got:%d, want:%d", len(ev.Code), len(wants))
	}

	for idx, got := range ev.Code {
		if got != wants[idx] {
			t.Errorf("Unexpected code at %d. got:0x%08x, want:0x%08x", idx, got, wants[idx])
		}
	}
}
