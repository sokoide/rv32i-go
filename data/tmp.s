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
	ret