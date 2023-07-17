package rv32i

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Loader struct {
	// Instructions *[]uint32
}

func NewLoader() *Loader {
	return &Loader{
		// Instructions: nil,
	}
}

func (p *Loader) LoadAt(filePath string, loadAddr []uint8, maxSize uint32) error {
	var err error
	var mem *[]uint32
	ext := filepath.Ext(filePath)

	if ext == ".txt" {
		mem, err = p.ReadText(filePath)
		for idx, u32 := range *mem {
			loadAddr[idx*4] = uint8(u32 & 0x000000FF)
			loadAddr[idx*4+1] = uint8((u32 & 0x0000FF00) >> 8)
			loadAddr[idx*4+2] = uint8((u32 & 0x00FF0000) >> 16)
			loadAddr[idx*4+3] = uint8((u32 & 0xFF000000) >> 24)
		}
	} else {
		// TODO
		mem, err = p.ReadBinary(filePath)
	}

	return err
}

func (p *Loader) ReadText(path string) (*[]uint32, error) {
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

func (p *Loader) ReadBinary(path string) (*[]uint32, error) {
	var ba []uint32

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return &ba, nil
}
