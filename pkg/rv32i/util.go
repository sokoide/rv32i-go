package rv32i

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func SignExtension(imm uint32, digit int) uint32 {
	sign := (imm >> digit) & 0b1
	mask := uint32(0xFFFFFFFF) - (uint32(1<<digit) - 1)
	if sign == 1 {
		imm = mask | imm
	}
	return imm
}

func ReadBinary(path string) (*[]uint32, error) {
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

func byteString2u8(b byte) uint8 {
	if b >= '0' && b <= '9' {
		return b - '0'
	} else if b >= 'a' && b <= 'f' {
		return b - 'a' + 10
	} else if b >= 'A' && b <= 'F' {
		return b - 'A' + 10
	} else {
		panic(fmt.Sprintf("byte 0x%02x not supported", b))
	}
}
