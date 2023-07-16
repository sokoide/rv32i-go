package rv32i

type Emulator struct {
	Cpu     *Cpu
	Program *Program
}

func NewEmulator() *Emulator {
	program := NewProgram()
	cpu := NewCpu(program)

	return &Emulator{
		Cpu:     cpu,
		Program: program,
	}
}

func (e *Emulator) Reset() {
	e.Cpu.Reset()
}

func (e *Emulator) Load(filePath string) error {
	return e.Program.Load(filePath)
}

func (e *Emulator) Step() error {
	return e.Cpu.Step()
}

func (e *Emulator) Dump() {
	e.Cpu.DumpRegisters()
}
