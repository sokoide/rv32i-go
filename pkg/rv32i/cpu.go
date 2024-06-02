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

var Regs = map[string]int{
	"zero": 0,
	"ra":   1,
	"sp":   2,
	"gp":   3,
	"tp":   4,
	"t0":   5,
	"t1":   6,
	"t2":   7,
	"s0":   8,
	"fp":   8,
	"s1":   9,
	"a0":   10,
	"a1":   11,
	"a2":   12,
	"a3":   13,
	"a4":   14,
	"a5":   15,
	"a6":   16,
	"a7":   17,
	"s2":   18,
	"s3":   19,
	"s4":   20,
	"s5":   21,
	"s6":   22,
	"s7":   23,
	"s8":   24,
	"s9":   25,
	"s10":  26,
	"s11":  27,
	"t3":   28,
	"t4":   29,
	"t5":   30,
	"t6":   31,
	"x0":   0,
	"x1":   1,
	"x2":   2,
	"x3":   3,
	"x4":   4,
	"x5":   5,
	"x6":   6,
	"x7":   7,
	"x8":   8,
	"x9":   9,
	"x10":  10,
	"x11":  11,
	"x12":  12,
	"x13":  13,
	"x14":  14,
	"x15":  15,
	"x16":  16,
	"x17":  17,
	"x18":  18,
	"x19":  19,
	"x20":  20,
	"x21":  21,
	"x22":  22,
	"x23":  23,
	"x24":  24,
	"x25":  25,
	"x26":  26,
	"x27":  27,
	"x28":  28,
	"x29":  29,
	"x30":  30,
	"x31":  31,
}

var RegsR = map[uint8]string{
	0:  "zero",
	1:  "ra",
	2:  "sp",
	3:  "gp",
	4:  "tp",
	5:  "t0",
	6:  "t1",
	7:  "t2",
	8:  "s0", // s0/fp
	9:  "s1",
	10: "a0",
	11: "a1", // a1/ret
	12: "a2", // a2/ret
	13: "a3",
	14: "a4",
	15: "a5",
	16: "a6",
	17: "a7",
	18: "s2",
	19: "s3",
	20: "s4",
	21: "s5",
	22: "s6",
	23: "s7",
	24: "s8",
	25: "s9",
	26: "s10",
	27: "s11",
	28: "t3",
	29: "t4",
	30: "t5",
	31: "t6",
}

func RegName(i uint8) string {
	return RegsR[i]
}

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
	trace("PC: 0x%08x, u32instr: %08x", c.PC, u32instr)

	// decode
	instr := NewInstruction(u32instr)
	trace("instr: %+v", instr)

	// execute
	incrementPC := c.Execute(instr)

	// increment PC if it's not jump
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
		trace("lui: X[%d] <- %x", i.Rd, i.Imm)
		if i.Rd > 0 {
			c.X[i.Rd] = i.Imm
		}
	case OpAuipc:
		trace("auipc: X[%d] <- PC:%x + imm:%x", i.Rd, c.PC, i.Imm)
		if i.Rd > 0 {
			c.X[i.Rd] = c.PC + i.Imm
		}
	case OpJal:
		t := c.PC + 4
		c.PC += i.Imm
		if i.Rd > 0 {
			c.X[i.Rd] = t
		}
		trace("jal: PC=%x, X[%d]=%x", c.PC, i.Rd, t)
		incrementPC = false
	case OpJalr:
		t := c.PC + 4
		c.PC = (c.X[i.Rs1] + i.Imm) & 0xffffffe
		if i.Rd > 0 {
			c.X[i.Rd] = t
		}
		trace("jalr: PC=%x, X[%d]=%x", c.PC, i.Rd, t)
		incrementPC = false
	case OpBeq:
		trace("beq: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] == c.X[i.Rs2] {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpBne:
		trace("bne: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] != c.X[i.Rs2] {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpBlt:
		// signed comparison
		a := int32(c.X[i.Rs1])
		b := int32(c.X[i.Rs2])
		trace("blt: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if a < b {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpBge:
		// signed comparison
		a := int32(c.X[i.Rs1])
		b := int32(c.X[i.Rs2])
		trace("bge: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if a >= b {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpBltu:
		// unsigned comparison
		trace("bltu: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] < c.X[i.Rs2] {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpBgeu:
		// unsigned comparison
		trace("bgeu: Rs1:%x, Rs2:%x", i.Rs1, i.Rs2)
		if c.X[i.Rs1] >= c.X[i.Rs2] {
			c.PC += i.Imm
			incrementPC = false
		}
	case OpLb:
		// sign extension
		addr := c.X[i.Rs1] + i.Imm
		trace("lb: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU8(addr))
		if i.Rd > 0 {
			c.X[i.Rd] = SignExtension(data, 7)
		}
	case OpLh:
		// sign extension
		addr := c.X[i.Rs1] + i.Imm
		trace("lh: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU16(addr))
		if i.Rd > 0 {
			c.X[i.Rd] = SignExtension(data, 15)
		}
	case OpLw:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		trace("lw: read %x -> X[%d]", addr, i.Rd)
		data := c.Emu.ReadU32(addr)
		if i.Rd > 0 {
			c.X[i.Rd] = data
		}
	case OpLbu:
		// zero extension
		addr := c.X[i.Rs1] + i.Imm
		trace("lbu: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU8(addr))
		if i.Rd > 0 {
			c.X[i.Rd] = data
		}
	case OpLhu:
		// zero extension
		addr := c.X[i.Rs1] + i.Imm
		trace("lhu: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU16(addr))
		if i.Rd > 0 {
			c.X[i.Rd] = data
		}
	case OpSb:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := uint8(c.X[i.Rs2] & 0xFF)
		trace("sb: write %x at %x", data, addr)
		c.Emu.WriteU8(addr, data)
	case OpSh:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := uint16(c.X[i.Rs2] & 0xFFFF)
		trace("sh: write %x at %x", data, addr)
		c.Emu.WriteU16(addr, data)
	case OpSw:
		// no extension
		addr := c.X[i.Rs1] + i.Imm
		data := c.X[i.Rs2]
		trace("sw: write %x at %x", data, addr)
		c.Emu.WriteU32(addr, data)
	case OpAddi:
		trace("addi: rs1:%x + imm:%x -> rd:%x", i.Rs1, i.Imm, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] + i.Imm
		}
	case OpSlti:
		// signed comparison
		trace("slti: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if int32(c.X[i.Rs1]) < int32(i.Imm) {
			if i.Rd > 0 {
				c.X[i.Rd] = 1
			}
		} else {
			c.X[i.Rd] = 0
		}
	case OpSltiu:
		// unsigned comparison
		trace("sltiu: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if c.X[i.Rs1] < i.Imm {
			if i.Rd > 0 {
				c.X[i.Rd] = 1
			}
		} else {
			c.X[i.Rd] = 0
		}
	case OpXori:
		trace("xori: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] ^ i.Imm
		}
	case OpOri:
		trace("ori: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] | i.Imm
		}
	case OpAndi:
		trace("andi: rs1:%x, imm:%x, rd:%x", i.Rs1, i.Imm, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] & i.Imm
		}
	case OpSlli:
		// logical shift
		trace("slli: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] << i.Rs2
		}
	case OpSrli:
		// logical shift
		trace("srli: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] >> i.Rs2
		}
	case OpSrai:
		// arithmetic shift
		trace("srai: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		data := c.X[i.Rs1] >> i.Rs2
		if i.Rd > 0 {
			c.X[i.Rd] = SignExtension(data, 31-int(i.Rs2))
		}
	case OpAdd:
		trace("add: rs1:%x + rs2:%x -> rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] + c.X[i.Rs2]
		}
	case OpSub:
		trace("sub: rs1:%x + rs2:%x -> rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] - c.X[i.Rs2]
		}
	case OpSll:
		// logical shift
		trace("sll: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] << c.X[i.Rs2]
		}
	case OpSlt:
		// signed comparison
		trace("slt: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if int32(c.X[i.Rs1]) < int32(c.X[i.Rs2]) {
			if i.Rd > 0 {
				c.X[i.Rd] = 1
			}
		} else {
			c.X[i.Rd] = 0
		}
	case OpSltu:
		// unsigned comparison
		trace("sltu: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if c.X[i.Rs1] < c.X[i.Rs2] {
			if i.Rd > 0 {
				c.X[i.Rd] = 1
			}
		} else {
			c.X[i.Rd] = 0
		}
	case OpXor:
		trace("xor: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] ^ c.X[i.Rs2]
		}
	case OpSrl:
		// logical shift
		trace("srl: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			shamt := 0b11111 & c.X[i.Rs2]
			c.X[i.Rd] = c.X[i.Rs1] >> shamt
		}
	case OpSra:
		// arithmetic shift
		trace("sra: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		shamt := 0b11111 & c.X[i.Rs2]
		data := c.X[i.Rs1] >> shamt
		if i.Rd > 0 {
			c.X[i.Rd] = SignExtension(data, 31-int(shamt))
		}
	case OpOr:
		trace("or: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] | c.X[i.Rs2]
		}
	case OpAnd:
		trace("and: rs1:%x, rs2:%x, rd:%x", i.Rs1, i.Rs2, i.Rd)
		if i.Rd > 0 {
			c.X[i.Rd] = c.X[i.Rs1] & c.X[i.Rs2]
		}
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
