package rv32i

const MaxMemory = uint32(0x10_000)

type Emulator struct {
	Cpu    *Cpu
	Memory []uint8
}

func NewEmulator() *Emulator {
	cpu := NewCpu()

	emu := Emulator{
		Cpu:    cpu,
		Memory: make([]uint8, MaxMemory),
	}
	cpu.Emu = &emu

	return &emu
}

func (e *Emulator) Reset() {
	e.Cpu.Reset()
	e.Memory = make([]uint8, MaxMemory)
}

func (e *Emulator) Load(filePath string) error {
	loader := NewLoader()
	return loader.LoadAt(filePath, &e.Memory, MaxMemory)
}

func (e *Emulator) LoadString(data string) error {
	loader := NewLoader()
	return loader.LoadStringAt(data, &e.Memory, MaxMemory)
}

func (e *Emulator) Step() error {
	return e.Cpu.Step()
}

func (e *Emulator) Run() error {
	var err error

	for {
		err = e.Cpu.Step()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Emulator) StepUntil(PC uint32) error {
	var err error
	for {
		if e.Cpu.PC == PC {
			break
		}
		err = e.Cpu.Step()
		if err != nil {
			return err
		}
	}

	return nil
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
	data = uint16(e.Memory[addr]) | uint16(e.Memory[addr+1])<<8

	return data
}

func (e *Emulator) ReadU32(addr uint32) uint32 {
	var data uint32
	data = uint32(e.Memory[addr]) | uint32(e.Memory[addr+1])<<8 | uint32(e.Memory[addr+2])<<16 | uint32(e.Memory[addr+3])<<24

	return data
}
