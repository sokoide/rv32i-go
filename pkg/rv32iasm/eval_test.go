package rv32iasm

import (
	"strings"
	"testing"
)

func Test_EvaluateProgram(t *testing.T) {
	src := `boot:
# This is a comment line
	li ra, 0 #0
	li s0, 0 #1 This is a comment
	lui a0, 4 #2
	auipc sp, 1 #3
	addi	sp, sp, -12 #4
	add	sp, sp, a0 #5
	jal riscv32_boot #6
_out:
	ret #7
is_even:
	addi	sp, sp, -16 #8
	sw	ra, 12(sp) #9
	sw	s0, 8(sp) #10
	addi	s0, sp, 16 #11
	sw	a0, -12(s0) #12
	lw	a0, -12(s0) #13
	srli	a1, a0, 31 #14
	add	a1, a0, a1 #15
	andi	a1, a1, -2
	sub	a0, a0, a1
	seqz a0, a0
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
		0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410113, 0x00a10133, 0x04c0006f,
		0x00008067,
		0xff010113, 0x00112623, 0x00812423, 0x01010413, 0xfea42a23, 0xff442503, 0x01f55593, 0x00b505b3,
		0xffe5f593, 0x40b50533, 0x00153513,
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
