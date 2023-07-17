package rv32i

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Register ABIName Description                         Saver
// -----------------------------------------------------------
// x0       zero    Hard-wired zero                     —
// x1       ra      Return address                      Caller
// x2       sp      Stack pointer                       Callee
// x3       gp      Global pointer                      —
// x4       tp      Thread pointer                      —
// x5–7     t0–2    Temporaries                         Caller
// x8       s0/fp   Saved register/frame pointer        Callee
// x9       s1      Saved register                      Callee
// x10–11   a0–1    Function arguments/return values    Caller
// x12–17   a2–7    Function arguments                  Caller
// x18–27   s2–11   Saved registers                     Callee
// x28–31   t3–6    Temporaries                         Caller

type Cpu struct {
	X   []uint32 // registers
	PC  uint32   // program counter
	Emu *Emulator
}

func NewCpu() *Cpu {
	return &Cpu{
		X:   make([]uint32, 32),
		PC:  0,
		Emu: nil,
	}
}

func (c *Cpu) Reset() {
	c.X = make([]uint32, 32)
	c.PC = 0
}

func (c *Cpu) Step() error {
	var err error
	var u32instr uint32

	// fetch
	u32instr, err = c.Fetch()
	if err != nil {
		return err
	}
	log.Tracef("PC: 0x%08x, u32instr: %08x", c.PC, u32instr)

	// decode
	instr := NewInstruction(u32instr)
	log.Tracef("instr: %+v", instr)

	// execute
	incrementPC := c.Execute(instr)

	// increment PC
	if incrementPC {
		c.PC += 4
	}

	return nil
}

func (c *Cpu) DumpRegisters() {
	log.Info("* Registers")
	for i := 0; i < len(c.X); i++ {
		log.Infof("x%d = %d, 0x%08x", i, c.X[i], c.X[i])
	}
	log.Infof("pc = 0x%08x", c.PC)
}

func (c *Cpu) Fetch() (uint32, error) {
	if c.PC > MaxMemory {
		return 0, errors.New("PC overflow")
	}
	i := c.Emu.ReadU32(c.PC)

	return i, nil
}

func (c *Cpu) Execute(i *Instruction) bool {
	var op OpName
	op = i.GetOpName()
	incrementPC := true

	switch op {
	case OpLui:
		log.Tracef("lui: X[%d] <- %x", i.Rd, i.Imm)
		c.X[i.Rd] = i.Imm
	case OpAuipc:
		log.Tracef("auipc: X[%d] <- PC:%x + imm:%x", i.Rd, c.PC, i.Imm)
		c.X[i.Rd] = c.PC + i.Imm
	case OpJal:
		t := c.PC + 4
		c.PC += i.Imm
		c.X[i.Rd] = t
		log.Tracef("jal: PC=%x, X[%d]=%x", c.PC, i.Rd, t)
		incrementPC = false
	case OpJalr:
		t := c.PC + 4
		c.PC = (c.X[i.Rs1] + i.Imm)
		c.X[i.Rd] = t
		log.Tracef("jalr: PC=%x, X[%d]=%x", c.PC, i.Rd, t)
		incrementPC = false
	case OpBeq:
		log.Tracef("beq: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] == c.X[i.Rs2] {
			c.PC += i.Imm
		}
	case OpBne:
		log.Tracef("bne: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] != c.X[i.Rs2] {
			c.PC += i.Imm
		}
	case OpBlt:
		// signed comparison
		a := int32(i.Rs1)
		b := int32(i.Rs2)
		log.Tracef("blt: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if a < b {
			c.PC += i.Imm
		}
	case OpBge:
		// signed comparison
		a := int32(i.Rs1)
		b := int32(i.Rs2)
		log.Tracef("bge: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if a >= b {
			c.PC += i.Imm
		}
	case OpBltu:
		// unsigned comparison
		log.Tracef("bltu: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] < c.X[i.Rs2] {
			c.PC += i.Imm
		}
	case OpBgeu:
		// unsigned comparison
		log.Tracef("bgeu: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] >= c.X[i.Rs2] {
			c.PC += i.Imm
		}
	case OpLb:
		// sign extension
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lb: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU8(addr))
		c.X[i.Rd] = SignExtension(data, 7)
	case OpLh:
		// sign extension
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lh: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU16(addr))
		c.X[i.Rd] = SignExtension(data, 15)
	case OpLw:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lw: read %x -> X[%d]", addr, i.Rd)
		data := c.Emu.ReadU32(addr)
		c.X[i.Rd] = data
	case OpLbu:
		// zero extension
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lbu: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU8(addr))
		c.X[i.Rd] = data
	case OpLhu:
		// zero extension
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lhu: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU16(addr))
		c.X[i.Rd] = data
	case OpSb:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := uint8(c.X[i.Rs2] & 0xFF)
		log.Tracef("sb: write %x at %x", data, addr)
		c.Emu.WriteU8(addr, data)
	case OpSh:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := uint16(c.X[i.Rs2] & 0xFFFF)
		log.Tracef("sh: write %x at %x", data, addr)
		c.Emu.WriteU16(addr, data)
	case OpSw:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := c.X[i.Rs2]
		log.Tracef("sw: write %x at %x", data, addr)
		c.Emu.WriteU32(addr, data)
	case OpAddi:
		log.Tracef("addi: rs1:%x + imm:%x -> rd:%x", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + i.Imm
	case OpSlti:
		// signed comparison
		log.Tracef("slti: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if int32(c.X[i.Rs1]) < int32(i.Imm) {
			c.X[i.Rd] = 1
		} else {
			c.X[i.Rd] = 0
		}
	case OpSltiu:
		// unsigned comparison
		log.Tracef("sltiu: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if c.X[i.Rs1] < i.Imm {
			c.X[i.Rd] = 1
		} else {
			c.X[i.Rd] = 0
		}
	case OpXori:
		log.Tracef("xori: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] ^ i.Imm
	case OpOri:
		log.Tracef("ori: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] | i.Imm
	case OpAndi:
		log.Tracef("andi: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] & i.Imm
	case OpSlli:
		// logical shift
		log.Tracef("slli: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] << i.Rs2
	case OpSrli:
		// logical shift
		log.Tracef("srli: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] >> i.Rs2
	case OpSrai:
		// arithmetic shift
		log.Tracef("srai: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		data := c.X[i.Rs1] >> i.Rs2
		c.X[i.Rd] = SignExtension(data, 31-int(i.Rs2))
	case OpAdd:
		log.Tracef("add: rs1:%x + rs2:%x -> rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + c.X[i.Rs2]
	case OpSub:
		log.Tracef("sub: rs1:%x + rs2:%x -> rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] - c.X[i.Rs2]
	case OpSll:
		// logical shift
		log.Tracef("sll: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] << c.X[i.Rs1]
	case OpSlt:
		// signed comparison
		log.Tracef("slt: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if int32(c.X[i.Rs1]) < int32(c.X[i.Rs2]) {
			c.X[i.Rd] = 1
		} else {
			c.X[i.Rd] = 0
		}
	case OpSltu:
		// unsigned comparison
		log.Tracef("sltu: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if c.X[i.Rs1] < c.X[i.Rs2] {
			c.X[i.Rd] = 1
		} else {
			c.X[i.Rd] = 0
		}
	case OpXor:
		log.Tracef("xor: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] ^ c.X[i.Rs2]
	case OpSrl:
		// logical shift
		log.Tracef("srl: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] >> c.X[i.Rs2]
	case OpSra:
		// arithmetic shift
		log.Tracef("sra: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		data := c.X[i.Rs1] >> c.X[i.Rs2]
		c.X[i.Rd] = SignExtension(data, 31-int(c.X[i.Rs2]))
	case OpOr:
		log.Tracef("or: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] | c.X[i.Rs2]
	case OpAnd:
		log.Tracef("and: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] & c.X[i.Rs2]
	case OpFence:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpFenceI:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpEcall:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpEbreak:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrw:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrs:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrc:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrwi:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrsi:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	case OpCsrrci:
		log.Warnf("Op %v is not implemented yet. rs1:%x, rs2:%x, rd:%x, imm:%x", op, i.Rs1, i.Rs2, i.Rd, i.Imm)
	default:
		panic(fmt.Sprintf("Op: %s invalid", op))
	}
	return incrementPC
}
