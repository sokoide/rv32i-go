package rv32i

import (
	"testing"
)

func Test_ReadBinary(t *testing.T) {
	prog := NewProgram()
	p, err := prog.ReadBinary("../../data/sample-binary-001.txt")
	if err != nil {
		t.Error(err)
	}

	if len(*p) != 3 {
		t.Errorf("Len %d not expected", len(*p))
	}

	for idx, want := range []uint32{0x00000093, 0x00000413, 0x00004537} {
		if (*p)[idx] != want {
			t.Errorf("got: 0x%08x, want: 0x%08x", (*p)[idx], want)
		}
	}
}
