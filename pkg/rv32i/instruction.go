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
	Type   InstructionType
	Imm    uint32
	Funct7 uint8
	Rs2    uint8
	Rs1    uint8
	Funct3 uint8
	Rd     uint8
	Opcode uint8
}

func NewInstruction(instr uint32) *Instruction {
	opcode := uint8(instr & 0b1111111)
	rd := uint8((instr >> 7) & 0b11111)
	funct3 := uint8((instr >> 12) & 0b111)
	rs1 := uint8((instr >> 15) & 0b11111)
	rs2 := uint8((instr >> 20) & 0b11111)
	funct7 := uint8(instr >> 25)
	imm := uint32(0)

	instance := Instruction{
		Imm:    imm,
		Funct7: funct7,
		Rs2:    rs2,
		Rs1:    rs1,
		Funct3: funct3,
		Rd:     rd,
		Opcode: opcode,
	}

	instance.Type = instance.GetInstructionType()
	switch instance.Type {
	case InstructionTypeU:
		imm = uint32(instr & 0b11111111_11111111_11110000_00000000)
		instance.Imm = imm
	case InstructionTypeJ:
		if opcode == 0b1101111 {
			// JAL
			imm20 := instr >> 31
			imm101 := instr >> 21 & 0b11_11111111
			imm11 := instr >> 20 & 0b1
			imm1912 := instr >> 12 & 0b11111111
			imm = imm20<<20 | imm101<<1 | imm11<<11 | imm1912<<12
			imm = SignExtension(imm, 20)
		} else {
			// JALR
			imm = instr >> 20
			imm = SignExtension(imm, 11)
		}
		instance.Imm = imm
	case InstructionTypeB:
		imm12 := instr >> 31
		imm1015 := instr >> 25 & 0b111111
		imm41 := instr >> 8 & 0b1111
		imm11 := instr >> 7 & 0b1
		imm = imm12<<12 | imm1015<<10 | imm41<<1 | imm11<<11
		imm = SignExtension(imm, 12)
	case InstructionTypeI:
		imm110 := instr >> 20
		imm = imm110
		imm = SignExtension(imm, 11)
	case InstructionTypeS:
		imm115 := instr >> 25
		imm40 := instr >> 7 & 0b11111
		imm = imm115<<5 | imm40
		imm = SignExtension(imm, 11)
	}

	return &instance
}

func (i *Instruction) GetInstructionType() InstructionType {
	switch i.Opcode {
	case 0b0110111, 0b0010111:
		return InstructionTypeU
	case 0b1101111, 0b1100111:
		return InstructionTypeJ
	case 0b1100011:
		return InstructionTypeB
	case 0b0000011:
		return InstructionTypeI
	case 0b0010011:
		if i.Funct3 == 0b001 || i.Funct3 == 0b101 {
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
		panic(fmt.Sprintf("Opcode 0x%07b ot Supported", i.Opcode))
	}
}
