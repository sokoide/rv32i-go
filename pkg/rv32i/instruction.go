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

//go:generate stringer -type OpName
type OpName int

const (
	OpLui OpName = iota
	OpAuipc
	OpJal
	OpJalr
	OpBeq
	OpBne
	OpBlt
	OpBge
	OpBltu
	OpBgeu
	OpLb
	OpLh
	OpLw
	OpLbu
	OpLhu
	OpSb
	OpSh
	OpSw
	OpAddi
	OpSlti
	OpSltiu
	OpXori
	OpOri
	OpAndi
	OpSlli
	OpSrli
	OpSrai
	OpAdd
	OpSub
	OpSll
	OpSlt
	OpSltu
	OpXor
	OpSrl
	OpSra
	OpOr
	OpAnd
	OpFence
	OpFenceI
	OpEcall
	OpEbreak
	OpCsrrw
	OpCsrrs
	OpCsrrc
	OpCsrrwi
	OpCsrrsi
	OpCsrrci
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
		imm105 := instr >> 25 & 0b111111
		imm41 := instr >> 8 & 0b1111
		imm11 := instr >> 7 & 0b1
		imm = imm12<<11 | imm105<<5 | imm41<<1 | imm11
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
	instance.Imm = imm

	return &instance
}

func GetCodeBase(t string, op1 int, op2 int, op3 int) uint32 {
	var code uint32
	switch t {
	case "J":
		imm20 := (uint32(op2) >> 20) & 0b1
		imm101 := (uint32(op2) >> 1) & 0b11_11111111
		imm11 := (uint32(op2) >> 11) & 0b1
		imm1912 := (uint32(op2) >> 12) & 0b11111111
		offset := imm20<<19 | imm101<<9 | imm11<<8 | imm1912
		code = offset<<12 | (uint32(op1) << 7) | 0b1101111
	case "B":
		imm12 := (uint32(op3) >> 12) & 0b1
		imm105 := (uint32(op3) >> 5) & 0b111111
		imm41 := (uint32(op3) >> 1) & 0b1111
		imm11 := (uint32(op3) >> 11) & 0b1
		code = imm12<<31 | imm105<<25 | (uint32(op2) << 20) | (uint32(op1) << 15) | (0b000 << 12) | (imm41 << 8) | (imm11 << 7) | 0b1100011
	default:
		panic(fmt.Sprintf("Type %s not supported", t))
	}
	return code
}

func GenCode(opn OpName, op1 int, op2 int, op3 int) uint32 {
	var code uint32
	switch opn {
	case OpLui:
		code = (uint32(op2) << 12) | (uint32(op1) << 7) | 0b0110111
		return code
	case OpAuipc:
		code = (uint32(op2) << 12) | (uint32(op1) << 7) | 0b0010111
		return code
	case OpJal:
		code = GetCodeBase("J", op1, op2, op3)
		return code
	case OpJalr:
		offset := (uint32(op2) & 0b1111_11111111)
		code = offset<<20 | (uint32(op3) << 15) | (uint32(op1) << 7) | 0b1100111
		return code
	case OpBeq:
		code = GetCodeBase("B", op1, op2, op3) | (0b000 << 12)
		return code
	case OpBne:
		code = GetCodeBase("B", op1, op2, op3) | (0b001 << 12)
		return code
	case OpBlt:
		code = GetCodeBase("B", op1, op2, op3) | (0b100 << 12)
		return code
	case OpBge:
		code = GetCodeBase("B", op1, op2, op3) | (0b101 << 12)
		return code
	case OpBltu:
		code = GetCodeBase("B", op1, op2, op3) | (0b110 << 12)
		return code
	case OpBgeu:
		code = GetCodeBase("B", op1, op2, op3) | (0b111 << 12)
		return code
	case OpAddi:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (uint32(op1) << 7) | 0b0010011
		return code
	case OpSlti:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (0b010 << 12) | (uint32(op1) << 7) | 0b0010011
		return code
	case OpSltiu:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (0b011 << 12) | (uint32(op1) << 7) | 0b0010011
		return code
	case OpAndi:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (0b111 << 12) | (uint32(op1) << 7) | 0b0010011
		return code
	case OpSrli:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (0b101 << 12) | (uint32(op1) << 7) | 0b0010011
		return code
	case OpAdd:
		code = (uint32(op3) << 20) | (uint32(op2) << 15) | (uint32(op1) << 7) | 0b0110011
		return code
	case OpSub:
		code = (0b01000 << 27) | (uint32(op3) << 20) | (uint32(op2) << 15) | (uint32(op1) << 7) | 0b0110011
		return code
	case OpLw:
		code := (uint32(op2) << 20) | (uint32(op3) << 15) | (0b010 << 12) | (uint32(op1) << 7) | 0b0000011
		return code
	case OpLbu:
		code := (uint32(op2) << 20) | (uint32(op3) << 15) | (0b100 << 12) | (uint32(op1) << 7) | 0b0000011
		return code
	case OpSb:
		imm115 := (uint32(op2) >> 5) & 0b1111111
		imm40 := uint32(op2) & 0b11111
		code := (imm115 << 25) | (uint32(op1) << 20) | (uint32(op3) << 15) | (0b000 << 12) | (imm40 << 7) | 0b0100011
		return code
	case OpSh:
		imm115 := (uint32(op2) >> 5) & 0b1111111
		imm40 := uint32(op2) & 0b11111
		code := (imm115 << 25) | (uint32(op1) << 20) | (uint32(op3) << 15) | (0b001 << 12) | (imm40 << 7) | 0b0100011
		return code
	case OpSw:
		imm115 := (uint32(op2) >> 5) & 0b1111111
		imm40 := uint32(op2) & 0b11111
		code := (imm115 << 25) | (uint32(op1) << 20) | (uint32(op3) << 15) | (0b010 << 12) | (imm40 << 7) | 0b0100011
		return code
	// TODO:
	default:
		return 1
	}
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
		panic(fmt.Sprintf("Opcode 0x%07b not Supported", i.Opcode))
	}
}

func (i *Instruction) GetOpName() OpName {
	switch i.Type {
	case InstructionTypeU:
		switch i.Opcode {
		case 0b0110111:
			return OpLui
		case 0b0010111:
			return OpAuipc
		default:
			panic(fmt.Sprintf("Opcode: %07b is invalid for %v", i.Opcode, i.Type))
		}
	case InstructionTypeJ:
		switch i.Opcode {
		case 0b1101111:
			return OpJal
		case 0b1100111:
			return OpJalr
		default:
			panic(fmt.Sprintf("Opcode: %07b is invalid for %v", i.Opcode, i.Type))
		}
	case InstructionTypeB:
		switch i.Funct3 {
		case 0b000:
			return OpBeq
		case 0b001:
			return OpBne
		case 0b100:
			return OpBlt
		case 0b101:
			return OpBge
		case 0b110:
			return OpBltu
		case 0b111:
			return OpBgeu
		default:
			panic(fmt.Sprintf("Funct3: %03b is invalid for %v", i.Funct3, i.Type))
		}
	case InstructionTypeI:
		switch i.Opcode {
		case 0b00000011:
			// L*
			switch i.Funct3 {
			case 0b000:
				return OpLb
			case 0b001:
				return OpLh
			case 0b010:
				return OpLw
			case 0b100:
				return OpLbu
			case 0b101:
				return OpLhu
			default:
				panic(fmt.Sprintf("Opcode: %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct3, i.Type))
			}
		case 0b0010011:
			// ADDI, SLTI, ...
			switch i.Funct3 {
			case 0b000:
				return OpAddi
			case 0b010:
				return OpSlti
			case 0b011:
				return OpSltiu
			case 0b100:
				return OpXori
			case 0b110:
				return OpOri
			case 0b111:
				return OpAndi
			default:
				panic(fmt.Sprintf("Opcode: %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct3, i.Type))
			}
		default:
			panic(fmt.Sprintf("Opcode: %07b is invalid for %v", i.Opcode, i.Type))
		}
	case InstructionTypeS:
		switch i.Funct3 {
		case 0b000:
			return OpSb
		case 0b001:
			return OpSh
		case 0b010:
			return OpSw
		default:
			panic(fmt.Sprintf("Funct3: %03b is invalid for %v", i.Funct3, i.Type))
		}
	case InstructionTypeR:
		switch i.Opcode {
		case 0b0010011:
			switch i.Funct3 {
			case 0b001:
				return OpSlli
			case 0b101:
				switch i.Funct7 {
				case 0b0000000:
					return OpSrli
				case 0b0100000:
					return OpSrli
				default:
					panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
				}
			default:
				panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
			}
		case 0b0110011:
			switch i.Funct3 {
			case 0b000:
				switch i.Funct7 {
				case 0b0000000:
					return OpAdd
				case 0b0100000:
					return OpSub
				default:
					panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
				}
			case 0b001:
				return OpSll
			case 0b010:
				return OpSlt
			case 0b011:
				return OpSltu
			case 0b100:
				return OpXor
			case 0b101:
				switch i.Funct7 {
				case 0b0000000:
					return OpSrl
				case 0b0100000:
					return OpSra
				default:
					panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
				}
			case 0b110:
				return OpOr
			case 0b111:
				return OpAnd
			default:
				panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
			}
		default:
			panic(fmt.Sprintf("Opcode: %07b, Funct7 %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct7, i.Funct3, i.Type))
		}
	case InstructionTypeF:
		switch i.Funct3 {
		case 0b000:
			return OpFence
		case 0b001:
			return OpFenceI
		default:
			panic(fmt.Sprintf("Opcode: %07b,  Funct3: %03b is invalid for %v", i.Opcode, i.Funct3, i.Type))
		}
	case InstructionTypeC:
		switch i.Funct3 {
		case 0b000:
			switch i.Rs2 {
			case 0:
				return OpEcall
			case 1:
				return OpEbreak
			default:
				panic(fmt.Sprintf("Opcode: %07b,  Rs2: %b, Funct3: %03b is invalid for %v", i.Opcode, i.Rs2, i.Funct3, i.Type))
			}

		case 0b001:
			return OpCsrrw
		case 0b010:
			return OpCsrrs
		case 0b011:
			return OpCsrrc
		case 0b101:
			return OpCsrrwi
		case 0b110:
			return OpCsrrsi
		case 0b111:
			return OpCsrrci
		default:
			panic(fmt.Sprintf("Opcode: %07b, Funct3: %03b is invalid for %v", i.Opcode, i.Funct3, i.Type))
		}
	default:
		panic(fmt.Sprintf("Opcode: %07b is invalid for %v", i.Opcode, i.Type))
	}
}
