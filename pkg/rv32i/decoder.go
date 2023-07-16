package rv32i

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(instr uint32) Instruction {
	return NewInstruction(instr)
}
