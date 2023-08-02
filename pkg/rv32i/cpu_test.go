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
	code = GenCode(OpBlt, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// lt - signed negative
	cpu.PC = 0x50
	cpu.X[3] = 0xffffffff
	cpu.X[4] = 0x1010
	code = GenCode(OpBlt, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not lt
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBlt, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Bge --------------------
	cpu.Reset()
	// ge
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBge, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not ge - signed negative
	cpu.PC = 0x50
	cpu.X[3] = 0xffffffff
	cpu.X[4] = 0x1000
	code = GenCode(OpBge, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not ge
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1010
	code = GenCode(OpBge, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Bltu --------------------
	cpu.Reset()
	// lt
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1010
	code = GenCode(OpBltu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not lt - unsigned
	cpu.PC = 0x50
	cpu.X[3] = 0xffffffff
	cpu.X[4] = 0x1010
	code = GenCode(OpBltu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not lt
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBltu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Bgeu --------------------
	cpu.Reset()
	// ge
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1000
	code = GenCode(OpBgeu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// ge - unsigned negative
	cpu.PC = 0x50
	cpu.X[3] = 0xffffffff
	cpu.X[4] = 0x1000
	code = GenCode(OpBgeu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != false {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50+1024 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// not ge
	cpu.PC = 0x50
	cpu.X[3] = 0x1000
	cpu.X[4] = 0x1010
	code = GenCode(OpBgeu, 3, 4, 1024)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if inc != true {
		t.Error("Wrong inc")
	}
	if cpu.PC != 0x50 {
		t.Errorf("Wrong PC %x", cpu.PC)
	}

	// Addi --------------------
	cpu.Reset()
	code = GenCode(OpAddi, 10, 0, 42)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 42 {
		t.Error("Wrong X10")
	}

	code = GenCode(OpAddi, 11, 10, 42)
	instr = NewInstruction(code)
	inc = cpu.Execute(instr)
	if cpu.X[11] != 84 {
		t.Error("Wrong X11")
	}

	code = GenCode(OpAddi, 11, 10, -41)
	instr = NewInstruction(code)
	inc = cpu.Execute(instr)
	if cpu.X[11] != 1 {
		t.Error("Wrong X11")
	}
	code = GenCode(OpAddi, 11, 10, -44)
	instr = NewInstruction(code)
	inc = cpu.Execute(instr)
	if cpu.X[11] != 0xfffffffe {
		t.Error("Wrong X11")
	}

	// Slti --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSlti, 3, 4, 101)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 1 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSlti, 3, 4, 100)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSlti, 3, 4, -1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Error("Wrong X3")
	}

	// Sltiu --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSltiu, 3, 4, 101)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 1 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSltiu, 3, 4, 100)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	// -1 is 0xffffffff in uint32
	code = GenCode(OpSltiu, 3, 4, -1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 1 {
		t.Error("Wrong X3")
	}

	// Xori --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm sign exteded to 0b11111111_11111111_11111111_00001111
	code = GenCode(OpXori, 3, 4, 0b1111_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00110011_00110011_00110011_11000011 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm 0b00000000_00000000_00000011_00001111
	code = GenCode(OpXori, 3, 4, 0b0011_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11001100_11001100_11001111_11000011 {
		t.Error("Wrong X3")
	}

	// Andi --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm sign exteded to 0b11111111_11111111_11111111_00001111
	code = GenCode(OpAndi, 3, 4, 0b1111_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11001100_11001100_11001100_00001100 {
		t.Error("Wrong X3")
	}

	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm 0b00000000_00000000_00000011_00001111
	code = GenCode(OpAndi, 3, 4, 0b0011_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00000000_00000000_00000000_00001100 {
		t.Error("Wrong X3")
	}

}
