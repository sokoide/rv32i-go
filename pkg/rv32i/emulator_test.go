package rv32i

import (
	"testing"
)

func Test_Step(t *testing.T) {
	var want uint32
	// Lui, Add used
	e := NewEmulator()
	e.Load("../../data/sample-binary-001.txt")

	want = 0
	if e.Cpu.X[10] != want {
		t.Errorf("x10 must be 0x%08x, but was 0x%08x", want, e.Cpu.X[10])
	}

	e.Step()
	e.Step()
	e.Step()

	want = 16 * 1024
	if e.Cpu.X[10] != want {
		t.Errorf("x10 must be 0x%08x, but was 0x%08x", want, e.Cpu.X[10])
	}
}

func Test_Step2(t *testing.T) {
	var want uint32
	// Auipc, Addi, Jal used
	e := NewEmulator()
	e.Load("../../data/sample-binary-002.txt")

	want = 0x00000000
	if e.Cpu.X[2] != want {
		t.Errorf("SP must be 0x%08x, but was 0x%08x", want, e.Cpu.X[2])
	}

	e.Step()
	want = 0x00000004
	if e.Cpu.PC != want {
		t.Errorf("PC must be 0x%08x, but was 0x%08x", want, e.Cpu.PC)
	}

	for i := 1; i < 7; i++ {
		e.Step()
	}

	want = 0x00000060
	if e.Cpu.PC != want {
		t.Errorf("PC must be 0x%08x, but was 0x%08x", want, e.Cpu.PC)
	}

	// 80000008: 37 45 00 00   lui     a0, 4 # a0 <- 16KB*4 == 0x4000
	// 8000000c <.Lpcrel_hi0>:
	// 8000000c: 17 11 00 00   auipc   sp, 1 # sp <- 1<<12 == 0x1000 + current PC (0xc) == 0x100c
	// 80000010: 13 01 41 ff   addi    sp, sp, -12 <- 0x0EE4 == 0x1000
	// 80000014: 33 01 a1 00   add     sp, sp, a0 <- -0x1000 + 0x4000
	want = 0x00005000
	if e.Cpu.X[2] != want {
		t.Errorf("SP must be 0x%08x, but was 0x%08x", want, e.Cpu.X[2])
	}
}

// func Test_Step3(t *testing.T) {
// 	// sw, lw, srli, sub, ...
// 	e := NewEmulator()
// 	e.Load("../../data/sample-binary-003.txt")
//
// 	for i := 0; i < 10; i++ {
// 		e.Step()
// 	}
//
// 	want := uint32(0x00000060)
// 	if e.Cpu.PC != want {
// 		t.Errorf("PC must be 0x%08x, but was 0x%08x", want, e.Cpu.PC)
// 	}
// }
