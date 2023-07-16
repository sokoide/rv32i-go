package rv32i

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Cpu struct {
	Regs    []uint32
	PC      int
	Program *Program
}

func NewCpu(p *Program) *Cpu {
	return &Cpu{
		Regs:    make([]uint32, 32),
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
	c.Execute(instr)

	// write back
	c.WriteBack()

	return nil
}

func (c *Cpu) DumpRegisters() {
	log.Info("* Registers")
	for i := 0; i < len(c.Regs); i++ {
		log.Infof("x%d = %d, 0x%08x", i, c.Regs[i], c.Regs[i])
	}
	log.Infof("pc = %d", c.PC)
}

func (c *Cpu) Fetch() (uint32, error) {
	if c.PC >= len(*c.Program.Instructions) {
		return 0, errors.New("PC overflow")
	}
	i := (*c.Program.Instructions)[c.PC]
	c.PC++

	return i, nil
}

func (c *Cpu) Execute(instr *Instruction) error {
	// TODO
	return nil
}

func (c *Cpu) WriteBack() error {
	// TODO
	return nil
}
