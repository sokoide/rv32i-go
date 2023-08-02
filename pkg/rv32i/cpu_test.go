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

	// Lb --------------------
	cpu.Reset()
	cpu.Emu.Memory[42] = 3
	cpu.Emu.Memory[43] = 1
	code = GenCode(OpLb, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 3 {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	cpu.Reset()
	cpu.X[1] = 100
	cpu.Emu.Memory[142] = 4
	cpu.Emu.Memory[143] = 1
	code = GenCode(OpLb, 10, 42, 1) // x10 <- 42(x1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 4 {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	cpu.Reset()
	cpu.Emu.Memory[42] = 0xff
	cpu.Emu.Memory[43] = 1
	code = GenCode(OpLb, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0xffffffff {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	// Lh --------------------
	cpu.Reset()
	cpu.Emu.Memory[42] = 3
	cpu.Emu.Memory[43] = 1
	cpu.Emu.Memory[44] = 1
	code = GenCode(OpLh, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0x0103 {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	cpu.Reset()
	cpu.Emu.Memory[42] = 0xff
	cpu.Emu.Memory[43] = 0xff
	code = GenCode(OpLh, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0xffffffff {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	// Lw --------------------
	cpu.Reset()
	cpu.X[1] = 100
	cpu.Emu.Memory[142] = 3
	cpu.Emu.Memory[143] = 1
	cpu.Emu.Memory[144] = 1
	cpu.Emu.Memory[145] = 1
	code = GenCode(OpLw, 10, 42, 1) // x10 <- 42(x1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0x01010103 {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	cpu.Reset()
	cpu.Emu.Memory[42] = 0xff
	cpu.Emu.Memory[43] = 0xff
	cpu.Emu.Memory[44] = 0xff
	cpu.Emu.Memory[45] = 0xff
	code = GenCode(OpLw, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0xffffffff {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	// Lbu --------------------
	cpu.Reset()
	cpu.Emu.Memory[42] = 0xff
	cpu.Emu.Memory[43] = 0
	code = GenCode(OpLbu, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0x000000ff {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}

	// Lhu --------------------
	cpu.Reset()
	cpu.Emu.Memory[42] = 0xff
	cpu.Emu.Memory[43] = 0xff
	cpu.Emu.Memory[44] = 0
	code = GenCode(OpLhu, 10, 42, 0) // x10 <- 42(x0)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[10] != 0x0000ffff {
		t.Errorf("Wrong X10 0x%08x", cpu.X[10])
	}
	// TODO: Sb --------------------
	// TODO: Sh --------------------
	// TODO: Sw --------------------

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
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSlti, 3, 4, 100)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSlti, 3, 4, -1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Sltiu --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSltiu, 3, 4, 101)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 1 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	code = GenCode(OpSltiu, 3, 4, 100)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 100
	// -1 is 0xffffffff in uint32
	code = GenCode(OpSltiu, 3, 4, -1)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 1 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
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
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm 0b00000000_00000000_00000011_00001111
	code = GenCode(OpXori, 3, 4, 0b0011_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11001100_11001100_11001111_11000011 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Ori --------------------
	cpu.Reset()
	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm sign exteded to 0b11111111_11111111_11111111_00001111
	code = GenCode(OpOri, 3, 4, 0b1111_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11111111_11111111_11111111_11001111 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm 0b00000000_00000000_00000011_00001111
	code = GenCode(OpOri, 3, 4, 0b0011_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11001100_11001100_11001111_11001111 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
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
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	cpu.X[3] = 0
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// imm 0b00000000_00000000_00000011_00001111
	code = GenCode(OpAndi, 3, 4, 0b0011_00001111)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00000000_00000000_00000000_00001100 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Slli --------------------
	cpu.Reset()
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// left logical shift 26 bits
	code = GenCode(OpSlli, 3, 4, 0b11010)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00110000_00000000_00000000_00000000 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Srli --------------------
	cpu.Reset()
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	// right logical shift 26 bits
	code = GenCode(OpSrli, 3, 4, 0b11010)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00000000_00000000_00000000_00110011 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Srai --------------------
	cpu.Reset()
	// right arithmetic shift 26 bits - negative
	cpu.X[4] = 0b11001100_11001100_11001100_11001100
	code = GenCode(OpSrai, 3, 4, 0b11010)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b11111111_11111111_11111111_11110011 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// right arithmetic shift 26 bits - positive
	cpu.X[4] = 0b01001100_11001100_11001100_11001100
	code = GenCode(OpSrai, 3, 4, 0b11010)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[3] != 0b00000000_00000000_00000000_00010011 {
		t.Errorf("Wrong X3, 0b%032b", cpu.X[3])
	}

	// Add --------------------
	cpu.Reset()
	cpu.X[3] = 123
	cpu.X[4] = 4567
	code = GenCode(OpAdd, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 4690 {
		t.Errorf("Wrong X1, %d", cpu.X[1])
	}

	// overflow
	cpu.X[3] = 0x0002
	cpu.X[4] = 0xffffffff
	code = GenCode(OpAdd, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 1 {
		t.Errorf("Wrong X1, %d", cpu.X[1])
	}

	// Sub --------------------
	cpu.Reset()
	cpu.X[3] = 123
	cpu.X[4] = 23
	code = GenCode(OpSub, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 100 {
		t.Errorf("Wrong X1, %d", cpu.X[1])
	}

	// underflow
	cpu.X[3] = 100
	cpu.X[4] = 102
	code = GenCode(OpSub, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0xfffffffe {
		t.Errorf("Wrong X1, %d", cpu.X[1])
	}

	// Sll --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 4
	code = GenCode(OpSll, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b11001100_11001100_11001100_11000000 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// Slt --------------------
	cpu.Reset()
	cpu.X[3] = 3
	cpu.X[4] = 4
	code = GenCode(OpSlt, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 1 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 3
	cpu.X[4] = 3
	code = GenCode(OpSlt, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0xfffffffe // -2
	cpu.X[4] = 3
	code = GenCode(OpSlt, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 1 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0xfffffffe // -2
	cpu.X[4] = 0xfffffffd // -3
	code = GenCode(OpSlt, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// Sltu --------------------
	cpu.Reset()
	cpu.X[3] = 3
	cpu.X[4] = 4
	code = GenCode(OpSltu, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 1 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 3
	cpu.X[4] = 3
	code = GenCode(OpSltu, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0xfffffffe // -2, but testing as uint
	cpu.X[4] = 3
	code = GenCode(OpSltu, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0xfffffffe // -2, but testing as uint
	cpu.X[4] = 0xfffffffd // -3, but testing as uint
	code = GenCode(OpSltu, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// Xor --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 0b11110000_11110000_11110000_11110000
	code = GenCode(OpXor, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b00111100_00111100_00111100_00111100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// Srl --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 8
	code = GenCode(OpSrl, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b00000000_11001100_11001100_11001100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	// only the lower 5bits are used -> 4 bit shift
	cpu.X[4] = 0b11111111_11111111_11111111_00000100
	code = GenCode(OpSrl, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b00001100_11001100_11001100_11001100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// Sra --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 8
	code = GenCode(OpSra, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b11111111_11001100_11001100_11001100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	cpu.X[3] = 0b10001100_11001100_11001100_11001100
	// only the lower 5bits are used -> 4 bit shift
	cpu.X[4] = 0b11111111_11111111_11111111_00000100
	code = GenCode(OpSra, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b11111000_11001100_11001100_11001100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}
	//  Or --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 0b11110000_11110000_11110000_11110000
	code = GenCode(OpOr, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b11111100_11111100_11111100_11111100 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}

	// And --------------------
	cpu.Reset()
	cpu.X[3] = 0b11001100_11001100_11001100_11001100
	cpu.X[4] = 0b11110000_11110000_11110000_11110000
	code = GenCode(OpAnd, 1, 3, 4)
	instr = NewInstruction(code)

	inc = cpu.Execute(instr)
	if cpu.X[1] != 0b11000000_11000000_11000000_11000000 {
		t.Errorf("Wrong X1, 0b%032b", cpu.X[1])
	}
}
