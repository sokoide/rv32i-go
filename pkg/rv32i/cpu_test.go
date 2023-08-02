package rv32i

import "testing"

func Test_Execute(t *testing.T) {
	var code uint32
	var instr *Instruction
	var inc bool

	cpu := NewCpu()
	emu := Emulator{
		Cpu:    cpu,
		Memory: make([]uint8, MaxMemory),
	}
	cpu.Emu = &emu

	// Lui --------------------
	cpu.Reset()
	code = GenCode(OpLui, 10, 41, 0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.X[10] != 41<<12 {
		t.Error("Wrong X10")
	}

	// Auipc --------------------
	cpu.Reset()
	// 0x50 + 41 << 12
	cpu.PC = 0x50
	code = GenCode(OpAuipc, 10, 41, 0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.X[10] != (41<<12)+0x50 {
		t.Error("Wrong X10")
	}

	// Jal --------------------
	cpu.Reset()
	// 0x50 + 1024
	cpu.PC = 0x50
	code = GenCode(OpJal, 11, 1024, 0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
	if cpu.X[11] != 0x50+4 {
		t.Error("Wrong X11")
	}
	// 0x50 -12
	cpu.PC = 0x50
	code = GenCode(OpJal, 11, -12, 0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50-12 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
	if cpu.X[11] != 0x50+4 {
		t.Error("Wrong X11")
	}

	cpu.Reset()
	// 0x50 -12, x0 not changed
	cpu.PC = 0x50
	code = GenCode(OpJal, 0, -12, 0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50-12 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
	if cpu.X[0] != 0 {
		t.Error("Wrong X0 (shold be always 0)")
	}

	// Jalr --------------------
	cpu.Reset()
	// 0x1000 + 1024
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	code = GenCode(OpJalr, 11, 1024, 3)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x1000+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
	if cpu.X[11] != 0x50+4 {
		t.Error("Wrong X11")
	}

	// 0x1000 -300
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	code = GenCode(OpJalr, 11, -300, 3)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x1000-300 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
	if cpu.X[11] != 0x50+4 {
		t.Error("Wrong X11")
	}

	// Beq --------------------
	cpu.Reset()
	// equal
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBeq, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not equal
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1002
	code = GenCode(OpBeq, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Bne --------------------
	cpu.Reset()
	// equal
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBne, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not equal
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1002
	code = GenCode(OpBne, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Blt --------------------
	cpu.Reset()
	// lt
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1010
	code = GenCode(OpBne, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	cpu.Reset()
	// not lt
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBne, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}
}
