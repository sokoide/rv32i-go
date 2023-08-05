package rv32iasm

import (
	"strings"
	"testing"

	"github.com/sokoide/rv32i-go/pkg/rv32i"
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
is_even:
	addi	sp, sp, -16
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	addi	s0, sp, 16
	sw	a0, -12(s0)
	lw	a0, -12(s0)
	srli	a1, a0, 31
	add	a1, a0, a1
	andi	a1, a1, -2
	sub	a0, a0, a1
	seqz a0, a0
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	addi	sp, sp, 16
	ret
riscv32_boot:
	addi	sp, sp, -16
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	addi	s0, sp, 16
	auipc	ra, 0
	jalr	24(ra)
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	addi	sp, sp, 16
	ret
main:
	addi	sp, sp, -32
	sw	ra, 28(sp)
	sw	s0, 24(sp)
	addi	s0, sp, 32
	li	a0, 10
	sw	a0, -12(s0)
	li	a0, 1
	sw	a0, -16(s0)
	lw	a0, -12(s0)
	auipc	ra, 0
	jalr	-136(ra)
	sb	a0, -17(s0)
	lw	a0, -16(s0)
	auipc	ra, 0
	jalr	-152(ra)
	sb	a0, -18(s0)
	lw	a0, -12(s0)
	auipc	ra, 0
	jalr	-172(ra)
	lw	a0, -16(s0)
	auipc	ra, 0
	jalr	-184(ra)
	lbu	a0, -17(s0)
	auipc	ra, 0
	jalr	-196(ra)
	lbu	a0, -18(s0)
	auipc	ra, 0
	jalr	-208(ra)
	li	a0, 0
	lw	ra, 28(sp)
	lw	s0, 24(sp)
	addi	sp, sp, 32
	ret
manualtest0:
	call manualtest1
	la t0, main
manualtest1:
	call main
	nop
	mv a1, a0
	neg a1, a0
	not a1, a0
	seqz a0, a1
	snez a0, a1
	sltz a0, a1
	sgtz a0, a1
	ret
`
	reader := strings.NewReader(src)
	scanner := NewScanner(reader)
	program, err := scanner.Parse()
	if err != nil {
		t.Error(err)
	}

	ev := NewEvaluator()
	_, err = ev.EvaluateProgram(program)
	if err != nil {
		t.Error("Failed to evaluate")
	}

	wants := []uint32{
		// boot ... TODO: replace the 6th code
		// 0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410113, 0x00a10133, 0x05c0006f,
		0x00000093, 0x00000413, 0x00004537, 0x00001117, 0xff410113, 0x00a10133, 0x044000ef,
		// _out
		0x00008067,
		// is_even
		0xff010113, 0x00112623, 0x00812423, 0x01010413, 0xfea42a23, 0xff442503, 0x01f55593, 0x00b505b3,
		0xffe5f593, 0x40b50533, 0x00153513, 0x00c12083, 0x00812403, 0x01010113, 0x00008067,
		// riscv32_boot
		0xff010113, 0x00112623, 0x00812423, 0x01010413, 0x00000097, 0x018080e7, 0x00c12083,
		0x00812403, 0x01010113, 0x00008067,
		// main
		0xfe010113, 0x00112e23, 0x00812c23, 0x02010413, 0x00a00513, 0xfea42a23, 0x00100513,
		0xfea42823, 0xff442503, 0x00000097, 0xf78080e7, 0xfea407a3, 0xff042503, 0x00000097,
		0xf68080e7, 0xfea40723, 0xff442503, 0x00000097, 0xf54080e7, 0xff042503, 0x00000097,
		0xf48080e7, 0xfef44503, 0x00000097, 0xf3c080e7, 0xfee44503, 0x00000097, 0xf30080e7,
		0x00000513, 0x01c12083, 0x01812403, 0x02010113, 0x00008067,
		// manualtest0 - call forward
		0x00000097, 0x118000e7, 0x00000297, 0x08400293,
		// manualtest1 - call backward
		0x00000097, 0x084000e7,
		// nop, mv
		0x00000013, 0x00050593,
		// neg, not
		0x40a005b3, 0xfff54593,
		// s*
		0x0015b513, 0x00b03533, 0x0005a533, 0x00b02533,
		// ret
		0x00008067,
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

func Test_Call_La(t *testing.T) {
	src := `entry:
	li x3,1
	li x4,2
	call hoge
	li a1, 42
	la t0, entry
	la t1, hoge
	ret
hoge:
	li a0, 123
	ret`
	// made by 'make tmp' in data dir
	//
	// src will be assembeld as below
	// 00000000 <entry>:
	//        0: 93 01 10 00   li      gp, 1
	//        4: 13 02 20 00   li      tp, 2
	//        8: 97 00 00 00   auipc   ra, 0
	//        c: e7 80 00 02   jalr    32(ra)
	//       10: 93 05 a0 02   li      a1, 42
	//       14: 97 02 00 00   auipc   t0, 0
	//       18: 93 82 c2 fe   addi    t0, t0, -20
	//       1c: 17 03 00 00   auipc   t1, 0
	//       20: 13 03 c3 00   addi    t1, t1, 12
	//       24: 67 80 00 00   ret
	//
	// 00000028 <hoge>:
	//       28: 13 05 b0 07   li      a0, 123
	//       2c: 67 80 00 00   ret
	want := uint32(123)

	reader := strings.NewReader(src)
	ev := NewEvaluator()
	code, err := ev.Assemble(reader)

	if err != nil {
		t.Error("Failed to assemble")
	}

	e := rv32i.NewEmulator()
	e.LoadString(strings.Join(code, "\n"))
	e.StepUntil(0x24)

	if e.Cpu.X[10] != want {
		t.Errorf("x10 must be 0x%08x, but was 0x%08x", want, e.Cpu.X[10])
	}

	want = uint32(0x0)
	if e.Cpu.X[5] != want {
		t.Errorf("x5 must be 0x%08x, but was 0x%08x", want, e.Cpu.X[5])
	}

	want = uint32(0x28)
	if e.Cpu.X[6] != want {
		t.Errorf("x6 must be 0x%08x, but was 0x%08x", want, e.Cpu.X[6])
	}
}
