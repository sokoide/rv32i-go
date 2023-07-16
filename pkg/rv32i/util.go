package rv32i

func SignExtension(imm uint32, digit int) uint32 {
	sign := (imm >> digit) & 0b1
	mask := uint32(0xFFFFFFFF) - (uint32(1<<digit) - 1)
	if sign == 1 {
		imm = mask | imm
	}
	return imm
}
