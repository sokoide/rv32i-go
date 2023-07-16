package rv32i

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type Program struct {
	Instructions *[]uint32
}

func NewProgram() *Program {
	return &Program{
		Instructions: nil,
	}
}

func (p *Program) Load(filePath string) error {
	var err error
	p.Instructions, err = p.ReadBinary(filePath)

	return err
}

func (p *Program) ReadBinary(path string) (*[]uint32, error) {
	var ba []uint32

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 1024)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return nil, errors.New("Line too long")
		}
		if err != nil {
			return nil, err
		}
		if len(line) < 21 {
			// not a vaild insruction line
			continue
		}
		if line[9] == '<' {
			// label
			continue
		}
		u32 := uint32(0)
		s := 0
		for i := 10; i < 21; i += 3 {
			by := byteString2u8(line[i])*16 + byteString2u8(line[i+1])
			u32 += uint32(by) << s
			s += 8
		}
		ba = append(ba, u32)
	}

	return &ba, nil
}
