package rv32i

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadAt(filePath string, loadAddr *[]uint8, maxSize uint32) error {
	var err error
	ext := filepath.Ext(filePath)

	if ext == ".txt" {
		fp, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fp.Close()

		reader := bufio.NewReaderSize(fp, 1024)
		err = l.readInto(reader, loadAddr, maxSize)
	} else {
		err = l.ReadBinary(filePath, loadAddr)
		if err != nil {
			return nil
		}
	}

	return err
}

func (l *Loader) LoadStringAt(data string, loadAddr *[]uint8, maxSize uint32) error {
	reader := bufio.NewReader(strings.NewReader(data))
	return l.readInto(reader, loadAddr, maxSize)
}

func (l *Loader) readInto(reader *bufio.Reader, loadAddr *[]uint8, maxSize uint32) error {
	var err error
	var mem *[]uint32

	mem, err = l.ReadText(reader)
	if err != nil {
		return nil
	}
	for idx, u32 := range *mem {
		(*loadAddr)[idx*4] = uint8(u32 & 0x000000FF)
		(*loadAddr)[idx*4+1] = uint8((u32 & 0x0000FF00) >> 8)
		(*loadAddr)[idx*4+2] = uint8((u32 & 0x00FF0000) >> 16)
		(*loadAddr)[idx*4+3] = uint8((u32 & 0xFF000000) >> 24)
	}

	return err
}

func (l *Loader) TextToBinary(pathIn string, pathOut string) error {
	var err error
	var p *[]uint32
	var fp *os.File

	fp, err = os.Open(pathIn)
	if err != nil {
		return err
	}
	defer fp.Close()

	p, err = l.ReadText(bufio.NewReaderSize(fp, 1024))

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

func (l *Loader) ReadText(reader *bufio.Reader) (*[]uint32, error) {
	var ba []uint32

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
		if len(line) < 20 {
			// not a vaild insruction line
			continue
		}
		if line[9] == '<' {
			// label
			continue
		}

		u32 := uint32(0)
		if string(line[10:12]) == "0x" {
			// syntax:
			//       0: 0x00000093 Addi ra, 0(zero)
			u64, _ := strconv.ParseUint(string(line[12:20]), 16, 32)
			u32 = uint32(u64)
		} else {
			// syntax:
			//       0: 93 00 00 00   li      ra, 0
			if len(line) < 21 {
				// not a vaild insruction line
				continue
			}
			s := 0
			for i := 10; i < 21; i += 3 {
				by := byteString2u8(line[i])*16 + byteString2u8(line[i+1])
				u32 += uint32(by) << s
				s += 8
			}
		}
		ba = append(ba, u32)
	}

	return &ba, nil
}

func (l *Loader) ReadBinary(path string, ba *[]byte) error {
	var err error
	var tmp []byte

	tmp, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	for idx, b := range tmp {
		(*ba)[idx] = b
	}

	return nil
}
