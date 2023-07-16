package rv32i

import "testing"

// *** registers ***
// x0: zero
// x1: ra return address ... caller saved
// x2: sp stack pointer ... callee saved
// x3: gp global pointer
// x4: tp thread pointer
// x5: t0 temporary/link register ... caller saved
// x6-x7: t1-t2 temporary .. caller saved
// x8: s0/fp save/frame pointer ... callee saved
// x9: s1 ... callee saved
// x10-x11: a0-a1 args/return value ... caller saved
// x12-x17: a2-a7 args ... caller saved
// x18-x27: s2-s11 save ... callee saved
// x28-x31: t3-t6 save ... caller saved

func Test_NewInstruction(t *testing.T) {
	type TestData struct {
		Instr uint32
		Want  Instruction
	}

	for _, td := range []TestData{
		// 80000088: 13 01 01 fe   addi    sp, sp, -32
		{0xfe010113, Instruction{Imm: 0, Funct7: 127, Rs2: 0, Rs1: 2, Funct3: 0, Rd: 2, Opcode: 0xfe010113 & 0b1111111}},
		// 8000008c: 23 2e 11 00   sw      ra, 28(sp)
		{0x00112e23, Instruction{Imm: 0, Funct7: 0, Rs2: 1, Rs1: 2, Funct3: 2, Rd: 28, Opcode: 0x00112e23 & 0b1111111}},
	} {
		got := NewInstruction(td.Instr)
		if got != td.Want {
			t.Errorf("NewInstruction failed 0x%x. got:%+v, want:%+v", td.Instr, got, td.Want)
		}
	}
}

func Test_GetInstructionType(t *testing.T) {
	type TestData struct {
		Instr uint32
		Want  InstructionType
	}

	for _, td := range []TestData{
		// 80000088: 13 01 01 fe   addi    sp, sp, -32
		{0xfe010113, InstructionTypeI},
		// 8000008c: 23 2e 11 00   sw      ra, 28(sp)
		{0x00112e23, InstructionTypeS},
		// 800000b0: e7 80 80 f7   jalr    -136(ra)
		{0xf78080e7, InstructionTypeJ},
		// 80000084: 67 80 00 00   ret
		{0x00008067, InstructionTypeJ},
		//       28: 63 00 00 00   beqz    zero, 0x28 <.Lline_table_start0+0x28>
		{0x00000063, InstructionTypeB},
		// 800000ac: 97 00 00 00   auipc   ra, 0
		{0x00000097, InstructionTypeU},
		//       3c: 73 63 76 31   csrrsi  t1, 791, 12
		{0x31766373, InstructionTypeC},
		//       9a: 73 2f 73 63   csrrs   t5, 1591, t1
		{0x63732f73, InstructionTypeC},
	} {
		opcode := uint8(td.Instr & 0b1111111)
		funct3 := uint8((td.Instr >> 12) & 0b111)
		got := GetInstructionType(opcode, funct3)

		if got != td.Want {
			t.Errorf("Decode failed for 0x%x. got:%v, want:%v", td.Instr, got, td.Want)
		}
	}
}

func Test_SignExtension(t *testing.T) {
	type TestData struct {
		Imm   uint32
		Digit int
		Want  uint32
	}

	for _, td := range []TestData{
		{0x0000FFFE, 15, 0xFFFFFFFE},
		{0x0000FFFF, 15, 0xFFFFFFFF},
		{0b00001000, 8, 0b00001000},
		{0b00001000, 3, 0b11111111_11111111_11111111_11111000},
	} {
		got := SignExtension(td.Imm, td.Digit)
		if got != td.Want {
			t.Errorf("SignExtention failed. in:%032b, got:%032b, want:%032b", td.Imm, got, td.Want)
		}
	}

}
