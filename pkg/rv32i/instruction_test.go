package rv32i

import "testing"

// *** rv32i registers ***
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
		{0xfe010113, Instruction{Type: InstructionTypeI, Imm: 4294967264, Funct7: 127, Rs2: 0, Rs1: 2, Funct3: 0, Rd: 2, Opcode: 0xfe010113 & 0b1111111}},
		// 8000008c: 23 2e 11 00   sw      ra, 28(sp)
		{0x00112e23, Instruction{Type: InstructionTypeS, Imm: 28, Funct7: 0, Rs2: 1, Rs1: 2, Funct3: 2, Rd: 28, Opcode: 0x00112e23 & 0b1111111}},
		// 800000b0: e7 80 80 f7   jalr    -136(ra)
		//  rs2 and funct7 are not used in JALR
		//  0xFFFFFEF0. It jumps to x[rs1] + 0xFFFFFEF0 == x0 + 0xFFFFFFE0 = -136
		{0xf78080e7, Instruction{Type: InstructionTypeJ, Imm: 0xFFFFFF78, Funct7: 123, Rs2: 24, Rs1: 1, Funct3: 0, Rd: 1, Opcode: 0xf78080e7 & 0b1111111}},
		// 80000010: ef 00 00 05  ▸jal▸0x80000060 <riscv32_boot>
		//  JAL only use rd and imm
		//  the current PC 0x80000010 + 0x50 = 0x80000060 is the jump target
		{0x050000ef, Instruction{Type: InstructionTypeJ, Imm: 0x50, Funct7: 2, Rs2: 16, Rs1: 0, Funct3: 0, Rd: 1, Opcode: 0x050000ef & 0b1111111}},
		// 80000084: 67 80 00 00   ret
		//  ret -> jalr zero, ra, 0
		{0x00008067, Instruction{Type: InstructionTypeJ, Imm: 0, Funct7: 0, Rs2: 0, Rs1: 1, Funct3: 0, Rd: 0, Opcode: 0x00008067 & 0b1111111}},
		//       28: 63 00 00 00   beqz    zero, 0x28 <.Lline_table_start0+0x28>
		{0x00000063, Instruction{Type: InstructionTypeB, Imm: 0, Funct7: 0, Rs2: 0, Rs1: 0, Funct3: 0, Rd: 0, Opcode: 0x00000063 & 0b1111111}},
		// 800000ac: 97 00 00 00   auipc   ra, 0
		{0x00000097, Instruction{Type: InstructionTypeU, Imm: 0, Funct7: 0, Rs2: 0, Rs1: 0, Funct3: 0, Rd: 1, Opcode: 0x00000097 & 0b1111111}},
		//       3c: 73 63 76 31   csrrsi  t1, 791, 12
		{0x31766373, Instruction{Type: InstructionTypeC, Imm: 0, Funct7: 24, Rs2: 23, Rs1: 12, Funct3: 6, Rd: 6, Opcode: 0x31766373 & 0b1111111}},
	} {
		got := NewInstruction(td.Instr)
		if *got != td.Want {
			t.Errorf("NewInstruction failed for 0x%x. got:%+v, want:%+v", td.Instr, got, td.Want)
		}
	}
}

func Test_GetOpName(t *testing.T) {
	type TestData struct {
		Instr uint32
		Want  OpName
	}

	for _, td := range []TestData{
		// 80000088: 13 01 01 fe   addi    sp, sp, -32
		{0xfe010113, OpAddi},
		// 8000008c: 23 2e 11 00   sw      ra, 28(sp)
		{0x00112e23, OpSw},
		// 800000b0: e7 80 80 f7   jalr    -136(ra)
		{0xf78080e7, OpJalr},
		// 80000010: ef 00 00 05  ▸jal▸0x80000060 <riscv32_boot>
		{0x050000ef, OpJal},
		// 80000084: 67 80 00 00   ret
		//  ret -> jalr zero, ra, 0
		{0x00008067, OpJalr},
		//       28: 63 00 00 00   beqz    zero, 0x28 <.Lline_table_start0+0x28>
		{0x00000063, OpBeq},
		// 800000ac: 97 00 00 00   auipc   ra, 0
		{0x00000097, OpAuipc},
		//       3c: 73 63 76 31   csrrsi  t1, 791, 12
		{0x31766373, OpCsrrsi},
	} {
		got := NewInstruction(td.Instr).GetOpName()
		if got != td.Want {
			t.Errorf("GetOpName failed for 0x%x. got:%+v, want:%+v", td.Instr, got, td.Want)
		}
	}
}

func Test_GenCode(t *testing.T) {
	type TestData struct {
		opn  OpName
		op1  int
		op2  int
		op3  int
		want uint32
	}

	tds := []TestData{
		{OpLui, 10, 4, 0, 0x00004537},
		{OpAuipc, 2, 1, 0, 0x00001117},
		{OpAddi, 8, 0, 0, 0x00000413},
	}

	for idx, td := range tds {
		got := GenCode(td.opn, td.op1, td.op2, td.op3)
		if got != td.want {
			t.Errorf("[%d] Wrong code. got: 0x%08x, want: 0x%08x", idx, got, td.want)
		}
	}
}
