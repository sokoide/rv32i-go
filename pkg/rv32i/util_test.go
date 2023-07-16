package rv32i

import (
	"testing"
)

func Test_SignExtension(t *testing.T) {
	type TestData struct {
		Imm   uint32
		Digit int
		Want  uint32
	}

	for _, td := range []TestData{
		{0x0000FFFE, 15, 0xFFFFFFFE},
		{0x0000FFFF, 15, 0xFFFFFFFF},
		{0b00001000, 8, 0b00001000},
		{0b00001000, 3, 0b11111111_11111111_11111111_11111000},
	} {
		got := SignExtension(td.Imm, td.Digit)
		if got != td.Want {
			t.Errorf("SignExtention failed. in:%032b, got:%032b, want:%032b", td.Imm, got, td.Want)
		}
	}

}
