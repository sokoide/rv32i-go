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
	auipc sp, 1
	addi	sp, sp, -12
	add	sp, sp, a0
	jal riscv32_boot
_out:
	ret
riscv32_boot:
`

	reader := strings.NewReader(src)
	scanner := NewScanner(reader)
	program, err := scanner.Parse()
	if err != nil {
		t.Error(err)
	}

	ev := NewEvaluator()
	err = ev.EvaluateProgram(program)
	if err != nil {
		t.Error("Failed to evaluate")
	}

	wants := []uint32{
		// TODO: after moving riscv32_boot to the right place
		// 0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410113, 0x00a10133, 0x044000ef, 0x00008067,
		0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410113, 0x00a10133, 0x0200006f, 0x00008067,
	}
	if len(ev.Code) != len(wants) {
		t.Errorf("Unexpected length. got:%d, want:%d", len(ev.Code), len(wants))
	}

	for idx, got := range ev.Code {
		if got != wants[idx] {
			t.Errorf("Unexpected code at %d. got:0x%08x, want:0x%08x", idx, got, wants[idx])
		}
	}
}
