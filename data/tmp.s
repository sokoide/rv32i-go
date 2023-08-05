# used by evalu_test.go
entry:
	li x3,1
	li x4,2
	call hoge
	li a1, 42
	la t0, entry
	la t1, hoge
	ret
hoge:
	li a0, 123
	ret
dummy:
	nop
	mv a1, a0
	neg a1, a0
	not a1, a0
	seqz a0, a1
	snez a0, a1
	sltz a0, a1
	sgtz a0, a1
	beqz a0, -4
	bnez a0, -4
	blez a0, -4
	bgez a0, -4
	bltz a0, -4
	bgtz a0, -4
	bgt a0, a1, -4
	ble a0, a1, -4
	bgtu a0, a1, -4
	bleu a0, a1, -4
	j 4
	jal 4
	jr a0
	jalr a0
	ret
