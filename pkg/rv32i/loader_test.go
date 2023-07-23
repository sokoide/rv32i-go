package rv32i

import (
	"bufio"
	"encoding/binary"
	"os"
	"testing"
)

func Test_ReadText(t *testing.T) {
	filePath := "../../data/sample-binary-001.txt"
	fp, err := os.Open(filePath)
	if err != nil {
		t.Errorf("failed to open %s", filePath)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 1024)

	loader := NewLoader()
	p, err := loader.ReadText(reader)
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

func Test_LoadStringAt(t *testing.T) {
	mem := make([]uint8, MaxMemory)
	program := `00000000 <boot>:
       0: 93 00 00 00   li      ra, 0
       4: 13 04 00 00   li      s0, 0
       8: 37 45 00 00   lui     a0, 4
`

	loader := NewLoader()
	err := loader.LoadStringAt(program, &mem, MaxMemory)
	if err != nil {
		t.Error("Failed to load a string")
	}

	testdata := []uint32{0x00000093, 0x00000413, 0x00004537}

	for idx, u32 := range testdata {
		got := binary.LittleEndian.Uint32(mem[idx*4 : idx*4+4])
		if got != u32 {
			t.Errorf("idx %d wanted 0x%08x, gt 0x%08x", idx, u32, got)
		}
	}
}
