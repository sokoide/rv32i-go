package rv32i

import (
	"fmt"
)

func SignExtension(imm uint32, digit int) uint32 {
	sign := (imm >> digit) & 0b1
	if sign == 1 {
		mask := uint32(0xFFFFFFFF) - (uint32(1<<digit) - 1)
		imm = mask | imm
	}
	return imm
}

func InterpretSingnedUint32(imm uint32) int {
	sign := (imm >> 31) & 0b1
	if sign == 1 {
		x := 0xffffffff ^ imm
		x++
		return int(x) * -1
	}
	return int(imm)
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

const intSize = 32 << (^uint(0) >> 63)

func Abs(v int) int {
	y := v >> (intSize - 1)
	return (v ^ y) - y
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}
