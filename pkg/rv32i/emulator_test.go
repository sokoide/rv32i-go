package rv32i

import "testing"

func Test_Step(t *testing.T) {
	// Lui, Add used
	e := NewEmulator()
	e.Load("../../data/sample-binary-001.txt")

	want := uint32(0)
	if e.Cpu.Regs[10] != want {
		t.Errorf("x10 must be 0x%08x, but was 0x%08x", want, e.Cpu.Regs[10])
	}

	e.Step()
	e.Step()
	e.Step()

	want = 16 * 1024
	if e.Cpu.Regs[10] != want {
		t.Errorf("x10 must be 0x%08x, but was 0x%08x", want, e.Cpu.Regs[10])
	}
}

func Test_Step2(t *testing.T) {
	// Auipc, Addi, Jal used
	e := NewEmulator()
	e.Load("../../data/sample-binary-002.txt")

	for i := 0; i < 7; i++ {
		e.Step()
	}

	want := uint32(0x00000060)
	if e.Cpu.PC != want {
		t.Errorf("PC must be 0x%08x, but was 0x%08x", want, e.Cpu.PC)
	}
}
