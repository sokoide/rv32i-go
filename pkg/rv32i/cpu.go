package rv32i

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Cpu struct {
	X       []uint32 // registers
	PC      uint32   // program counter
	Program *Program
	Emu     *Emulator
}

func NewCpu(p *Program) *Cpu {
	return &Cpu{
		X:       make([]uint32, 32),
		PC:      0,
		Program: p,
		Emu:     nil,
	}
}

func (c *Cpu) Reset() {
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
	log.Tracef("u32instr: %08x", u32instr)

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
	if c.PC/4 >= uint32(len(*c.Program.Instructions)) {
		return 0, errors.New("PC overflow")
	}
	i := (*c.Program.Instructions)[c.PC/4]

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
		panic("Not implemented yet")
	case OpBne:
		panic("Not implemented yet")
	case OpBlt:
		panic("Not implemented yet")
	case OpBge:
		panic("Not implemented yet")
	case OpBltu:
		panic("Not implemented yet")
	case OpBgeu:
		panic("Not implemented yet")
	case OpLb:
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lb: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU8(addr))
		c.X[i.Rd] &= 0xFFFFFF00
		c.X[i.Rd] |= data
	case OpLh:
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lh: read %x -> X[%d]", addr, i.Rd)
		data := uint32(c.Emu.ReadU16(addr))
		c.X[i.Rd] &= 0xFFFF0000
		c.X[i.Rd] |= data
	case OpLw:
		addr := c.X[i.Rs1] + i.Imm
		log.Tracef("lw: read %x -> X[%d]", addr, i.Rd)
		data := c.Emu.ReadU32(addr)
		c.X[i.Rd] = data
	case OpLbu:
		panic("Not implemented yet")
	case OpLhu:
		panic("Not implemented yet")
	case OpSb:
		addr := c.X[i.Rs1] + i.Imm
		data := uint8(c.X[i.Rs2] & 0xFF)
		log.Tracef("sb: write %x at %x", data, addr)
		c.Emu.WriteU8(addr, data)
	case OpSh:
		addr := c.X[i.Rs1] + i.Imm
		data := uint16(c.X[i.Rs2] & 0xFFFF)
		log.Tracef("sh: write %x at %x", data, addr)
		c.Emu.WriteU16(addr, data)
	case OpSw:
		addr := c.X[i.Rs1] + i.Imm
		data := c.X[i.Rs2]
		log.Tracef("sw: write %x at %x", data, addr)
		c.Emu.WriteU32(addr, data)
	case OpAddi:
		log.Tracef("addi: rs1:%x + imm:%x -> rd:%x", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + i.Imm
	case OpSlti:
		panic("Not implemented yet")
	case OpSltiu:
		panic("Not implemented yet")
	case OpXori:
		panic("Not implemented yet")
	case OpOri:
		panic("Not implemented yet")
	case OpAndi:
		panic("Not implemented yet")
	case OpSlli:
		panic("Not implemented yet")
	case OpSrli:
		panic("Not implemented yet")
	case OpSrai:
		panic("Not implemented yet")
	case OpAdd:
		log.Tracef("add: rs1:%x + rs2:%x -> rd:%x", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + c.X[i.Rs2]
	case OpSub:
		panic("Not implemented yet")
	case OpSll:
		panic("Not implemented yet")
	case OpSlt:
		panic("Not implemented yet")
	case OpSltu:
		panic("Not implemented yet")
	case OpXor:
		panic("Not implemented yet")
	case OpSrl:
		panic("Not implemented yet")
	case OpSra:
		panic("Not implemented yet")
	case OpOr:
		panic("Not implemented yet")
	case OpAnd:
		panic("Not implemented yet")
	case OpFence:
		panic("Not implemented yet")
	case OpFenceI:
		panic("Not implemented yet")
	case OpEcall:
		panic("Not implemented yet")
	case OpEbreak:
		panic("Not implemented yet")
	case OpCsrrw:
		panic("Not implemented yet")
	case OpCsrrs:
		panic("Not implemented yet")
	case OpCsrrc:
		panic("Not implemented yet")
	case OpCsrrwi:
		panic("Not implemented yet")
	case OpCsrrsi:
		panic("Not implemented yet")
	case OpCsrrci:
		panic("Not implemented yet")
	default:
		panic(fmt.Sprintf("Op: %s invalid", op))
	}
	return incrementPC
}
