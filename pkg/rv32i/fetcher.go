package rv32i

import "errors"

type Fetcher struct {
	PC           int
	Instructions *[]uint32
}

type RawInstruction struct {
	Instr uint32
}

func NewFetcher(instructions *[]uint32) *Fetcher {
	return &Fetcher{
		PC:           0,
		Instructions: instructions,
	}
}

func (f *Fetcher) Fetch() (uint32, error) {
	if f.PC >= len(*f.Instructions) {
		return 0, errors.New("PC overflow")
	}
	i := (*f.Instructions)[f.PC]
	f.PC++

	return i, nil
}
