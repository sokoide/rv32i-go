package rv32i

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadAt(filePath string, loadAddr []uint8, maxSize uint32) error {
	var err error
	var mem *[]uint32
	ext := filepath.Ext(filePath)

	if ext == ".txt" {
		mem, err = l.ReadText(filePath)
	} else {
		mem, err = l.ReadBinary(filePath)
	}

	for idx, u32 := range *mem {
		loadAddr[idx*4] = uint8(u32 & 0x000000FF)
		loadAddr[idx*4+1] = uint8((u32 & 0x0000FF00) >> 8)
		loadAddr[idx*4+2] = uint8((u32 & 0x00FF0000) >> 16)
		loadAddr[idx*4+3] = uint8((u32 & 0xFF000000) >> 24)
	}

	return err
}

func (l *Loader) TextToBinary(pathIn string, pathOut string) error {
	var err error
	var p *[]uint32
	var fp *os.File

	p, err = l.ReadText(pathIn)

	if err != nil {
		return nil
	}

	fp, err = os.OpenFile(pathOut, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, u32 := range *p {
		by4 := make([]byte, 4)
		by4[0] = byte(u32 & 0x000000FF)
		by4[1] = byte((u32 & 0x0000FF00) >> 8)
		by4[2] = byte((u32 & 0x00FF0000) >> 16)
		by4[3] = byte((u32 & 0xFF000000) >> 24)
		fp.Write(by4)
	}
	return nil
}

func (l *Loader) ReadText(path string) (*[]uint32, error) {
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

func (l *Loader) ReadBinary(path string) (*[]uint32, error) {
	var ba []uint32

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	return &ba, nil
}
