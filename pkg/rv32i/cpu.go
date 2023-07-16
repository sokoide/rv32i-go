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
}

func NewCpu(p *Program) *Cpu {
	return &Cpu{
		X:       make([]uint32, 32),
		PC:      0,
		Program: p,
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
		c.X[i.Rd] = i.Imm
	case OpAuipc:
		c.X[i.Rd] = c.PC + i.Imm
	case OpAddi:
		log.Tracef("addi: rs1:%x + imm:%x -> rd:%x", i.Rs1, i.Imm, i.Rd)
		fmt.Printf("addi: rs1:%x + imm:%x -> rd:%x\n", i.Rs1, i.Imm, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + i.Imm
	case OpAdd:
		fmt.Printf("add: rs1:%x + rs2:%x -> rd:%x\n", i.Rs1, i.Rs2, i.Rd)
		c.X[i.Rd] = c.X[i.Rs1] + c.X[i.Rs2]
	case OpSw:
		addrTarget := c.X[i.Rs1] + i.Imm
		data := c.X[i.Rs2]
		log.Tracef("writing %x at %08x", data, addrTarget)
	case OpJal:
		t := c.PC + 4
		c.PC += i.Imm
		c.X[i.Rd] = t
		incrementPC = false
	case OpJalr:
		t := c.PC + 4
		c.PC = (c.X[i.Rs1] + i.Imm)
		c.X[i.Rd] = t
		incrementPC = false
	default:
		// TODO: must implemente all operators
		panic(fmt.Sprintf("Op: %s not supproted yet", op))
	}
	return incrementPC
}
