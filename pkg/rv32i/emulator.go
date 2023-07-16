package rv32i

type Emulator struct {
	Cpu     *Cpu
	Memory  []uint8
	Program *Program
}

func NewEmulator() *Emulator {
	program := NewProgram()
	cpu := NewCpu(program)

	emu := Emulator{
		Cpu:     cpu,
		Program: program,
		Memory:  make([]uint8, 0x10_000),
	}
	cpu.Emu = &emu

	return &emu
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

func (e *Emulator) WriteU8(addr uint32, data uint8) {
	e.Memory[addr] = data
}

func (e *Emulator) WriteU16(addr uint32, data uint16) {
	e.Memory[addr] = uint8(data & 0x00FF)
	addr++
	e.Memory[addr] = uint8((data & 0xFF00) >> 8)
	addr++
}

func (e *Emulator) WriteU32(addr uint32, data uint32) {
	e.Memory[addr] = uint8(data & 0x000000FF)
	addr++
	e.Memory[addr] = uint8((data & 0x0000FF00) >> 8)
	addr++
	e.Memory[addr] = uint8((data & 0x00FF0000) >> 16)
	addr++
	e.Memory[addr] = uint8((data & 0xFF000000) >> 24)
}

func (e *Emulator) ReadU8(addr uint32) uint8 {
	return e.Memory[addr]
}

func (e *Emulator) ReadU16(addr uint32) uint16 {
	var data uint16
	data = uint16(e.Memory[addr]) | uint16(e.Memory[addr+1]<<8)

	return data
}

func (e *Emulator) ReadU32(addr uint32) uint32 {
	var data uint32
	data = uint32(e.Memory[addr]) | uint32(e.Memory[addr+1]<<8) | uint32(e.Memory[addr+2]<<14) | uint32(e.Memory[addr+3]<<24)

	return data
}
