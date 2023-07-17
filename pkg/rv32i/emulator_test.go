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

func Test_Step3(t *testing.T) {
	var want uint32
	var got uint32
	var wantu8 uint8
	var gotu8 uint8

	// sw, lw, srli, sub, ...
	e := NewEmulator()
	e.Load("../../data/sample-binary-003.txt")

	e.StepUntil(0x70)
	// now at 0x00000070
	got = e.Cpu.X[2]
	want = uint32(0x00004ff0)
	if got != want {
		t.Errorf("SP must be 0x%08x, but was 0x%08x", want, got)
	}

	// 5c: 13 01 01 ff  ▸addi▸   sp, sp, -16 -> sp == 0x4ff0
	// 60: 23 26 11 00  ▸sw▸ ra, 12(sp) -> stored ra at 0x4ffc
	// 'ra' must be the return address of the callee (<boot>)
	// which is 0x1c (next IP of 18: ef 00 40 04  ▸jal▸0x5c <riscv32_boot>
	wantu8 = uint8(0x1c)
	gotu8 = e.Memory[0x4ffc]
	if gotu8 != wantu8 {
		t.Errorf("0x4ffc must be 0x%02x, but was 0x%02x", wantu8, gotu8)
	}

	e.StepUntil(0x24)
	// a0 (1st argument to is_even
	got = e.Cpu.X[10]
	want = uint32(0x0000000a)
	if got != want {
		t.Errorf("a0 must be 0x%08x, was 0x%08x", want, got)
	}

	e.StepUntil(0x3c)
	// now at is_even 0x0000003c
	got = e.Cpu.X[11]
	want = uint32(0x00000000)
	if got != want {
		t.Errorf("a1 must be 0x%08x, but was 0x%08x", want, e.Cpu.PC)
	}

	e.StepUntil(0x58)
	e.Step()
	// 'ret' must return to the main function at 0xb0
	got = e.Cpu.PC
	want = uint32(0x000000b0)
	if got != want {
		t.Errorf("PC must be 0x%08x, but was 0x%08x", want, got)
	}
	// a0 == is_even(10) must be 1
	got = e.Cpu.X[10]
	want = uint32(0x0000001)
	if got != want {
		t.Errorf("a0 must be 0x%08x, but was 0x%08x", want, got)
	}

	e.StepUntil(0xb4)
	got = e.Cpu.X[10]
	want = uint32(0x0000001)
	if got != want {
		t.Errorf("a0 must be 0x%08x, but was 0x%08x", want, got)
	}

	e.StepUntil(0xc0)
	// a0 == is_even(1) must be 0
	got = e.Cpu.X[10]
	want = uint32(0x0000000)
	if got != want {
		t.Errorf("a0 must be 0x%08x, but was 0x%08x", want, got)
	}

	e.StepUntil(0x104)
}
