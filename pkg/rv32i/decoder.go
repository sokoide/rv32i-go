package rv32i

import "fmt"

//go:generate stringer -type InstructionType
type InstructionType int

const (
	InstructionTypeU InstructionType = iota
	InstructionTypeJ
	InstructionTypeB
	InstructionTypeI
	InstructionTypeS
	InstructionTypeR
	InstructionTypeF
	InstructionTypeC
)

type Instruction struct {
	Imm    uint32
	Funct7 uint8
	Rs2    uint8
	Rs1    uint8
	Funct3 uint8
	Rd     uint8
	Opcode uint8
}

func NewInstruction(instr uint32) Instruction {
	opcode := uint8(instr & 0b1111111)
	rd := uint8((instr >> 7) & 0b11111)
	funct3 := uint8((instr >> 12) & 0b111)
	rs1 := uint8((instr >> 15) & 0b11111)
	rs2 := uint8((instr >> 20) & 0b11111)
	funct7 := uint8(instr >> 25)
	imm := uint32(0)

	switch GetInstructionType(opcode, funct3) {
	case InstructionTypeU:
		imm = uint32(instr & 0b11111111_11111111_11110000_00000000)
	case InstructionTypeJ:
		imm20 := uint32(instr >> 31)
		imm101 := uint32((instr >> 21) & 0b1_11111111)
		imm11 := uint32((instr >> 20) & 0b1)
		imm1912 := uint32((instr >> 12) & 0b11111111)
		imm = (imm20 << 20) | (imm101 << 1) | (imm11 << 11) | (imm1912 << 12)
		imm = SignExtension(imm, 20)
	}

	return Instruction{
		Imm:    imm,
		Funct7: funct7,
		Rs2:    rs2,
		Rs1:    rs1,
		Funct3: funct3,
		Rd:     rd,
		Opcode: opcode,
	}
}

func GetInstructionType(opcode uint8, funct3 uint8) InstructionType {
	switch opcode {
	case 0b0110111, 0b0010111:
		return InstructionTypeU
	case 0b1101111, 0b1100111:
		return InstructionTypeJ
	case 0b1100011:
		return InstructionTypeB
	case 0b0000011:
		return InstructionTypeI
	case 0b0010011:
		if funct3 == 0b001 || funct3 == 0b101 {
			return InstructionTypeR
		}
		return InstructionTypeI
	case 0b0100011:
		return InstructionTypeS
	case 0b0110011:
		return InstructionTypeR
	case 0b0001111:
		return InstructionTypeF
	case 0b1110011:
		return InstructionTypeC
	default:
		panic(fmt.Sprintf("Opcode 0x%07b ot Supported", opcode))
	}
}

func SignExtension(imm uint32, digit int) uint32 {
	sign := (imm >> digit) & 0b1
	mask := uint32(0xFFFFFFFF) - (uint32(1<<digit) - 1)
	if sign == 1 {
		imm = mask | imm
	}
	return imm
}

func (i *Instruction) InstructionType() InstructionType {
	return InstructionTypeC
}
