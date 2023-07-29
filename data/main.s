	.text
	.attribute	4, 16
	.attribute	5, "rv32i2p0"
	.file	"main.c"
	.globl	is_even
	.p2align	2
	.type	is_even,@function
is_even:
.Lfunc_begin0:
	.file	0 "/Users/scott/repo/sokoide/rv32i-go/data" "main.c" md5 0xd9eb863840895378e9a64fab62ae38d5
	.loc	0 4 0
	.cfi_sections .debug_frame
	.cfi_startproc
	addi	sp, sp, -16
	.cfi_def_cfa_offset 16
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	.cfi_offset ra, -4
	.cfi_offset s0, -8
	addi	s0, sp, 16
	.cfi_def_cfa s0, 0
	sw	a0, -12(s0)
.Ltmp0:
	.loc	0 4 30 prologue_end
	lw	a0, -12(s0)
	.loc	0 4 32 is_stmt 0
	srli	a1, a0, 31
	add	a1, a0, a1
	andi	a1, a1, -2
	sub	a0, a0, a1
	.loc	0 4 36
	seqz	a0, a0
	.loc	0 4 23
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	.loc	0 4 23 epilogue_begin
	addi	sp, sp, 16
	ret
.Ltmp1:
.Lfunc_end0:
	.size	is_even, .Lfunc_end0-is_even
	.cfi_endproc

	.globl	riscv32_boot
	.p2align	2
	.type	riscv32_boot,@function
riscv32_boot:
.Lfunc_begin1:
	.loc	0 7 0 is_stmt 1
	.cfi_startproc
	addi	sp, sp, -16
	.cfi_def_cfa_offset 16
.Ltmp2:
	.loc	0 12 5 prologue_end
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	.cfi_offset ra, -4
	.cfi_offset s0, -8
	addi	s0, sp, 16
	.cfi_def_cfa s0, 0
	call	main
	.loc	0 13 1
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	.loc	0 13 1 epilogue_begin is_stmt 0
	addi	sp, sp, 16
	ret
.Ltmp3:
.Lfunc_end1:
	.size	riscv32_boot, .Lfunc_end1-riscv32_boot
	.cfi_endproc

	.globl	main
	.p2align	2
	.type	main,@function
main:
.Lfunc_begin2:
	.loc	0 15 0 is_stmt 1
	.cfi_startproc
	addi	sp, sp, -32
	.cfi_def_cfa_offset 32
	sw	ra, 28(sp)
	sw	s0, 24(sp)
	.cfi_offset ra, -4
	.cfi_offset s0, -8
	addi	s0, sp, 32
	.cfi_def_cfa s0, 0
	li	a0, 10
.Ltmp4:
	.loc	0 18 7 prologue_end
	sw	a0, -12(s0)
	li	a0, 1
	.loc	0 19 7
	sw	a0, -16(s0)
	.loc	0 20 17
	lw	a0, -12(s0)
	.loc	0 20 9 is_stmt 0
	call	is_even
	.loc	0 20 7
	sb	a0, -17(s0)
	.loc	0 21 17 is_stmt 1
	lw	a0, -16(s0)
	.loc	0 21 9 is_stmt 0
	call	is_even
	.loc	0 21 7
	sb	a0, -18(s0)
	.loc	0 22 10 is_stmt 1
	lw	a0, -12(s0)
	.loc	0 22 5 is_stmt 0
	call	_out
	.loc	0 23 10 is_stmt 1
	lw	a0, -16(s0)
	.loc	0 23 5 is_stmt 0
	call	_out
	.loc	0 24 10 is_stmt 1
	lbu	a0, -17(s0)
	.loc	0 24 5 is_stmt 0
	call	_out
	.loc	0 25 10 is_stmt 1
	lbu	a0, -18(s0)
	.loc	0 25 5 is_stmt 0
	call	_out
	li	a0, 0
	.loc	0 26 5 is_stmt 1
	lw	ra, 28(sp)
	lw	s0, 24(sp)
	.loc	0 26 5 epilogue_begin is_stmt 0
	addi	sp, sp, 32
	ret
.Ltmp5:
.Lfunc_end2:
	.size	main, .Lfunc_end2-main
	.cfi_endproc

	.section	.debug_abbrev,"",@progbits
	.byte	1
	.byte	17
	.byte	1
	.byte	37
	.byte	37
	.byte	19
	.byte	5
	.byte	3
	.byte	37
	.byte	114
	.byte	23
	.byte	16
	.byte	23
	.byte	27
	.byte	37
	.byte	17
	.byte	27
	.byte	18
	.byte	6
	.byte	115
	.byte	23
	.byte	0
	.byte	0
	.byte	2
	.byte	46
	.byte	1
	.byte	17
	.byte	27
	.byte	18
	.byte	6
	.byte	64
	.byte	24
	.byte	3
	.byte	37
	.byte	58
	.byte	11
	.byte	59
	.byte	11
	.byte	39
	.byte	25
	.byte	73
	.byte	19
	.byte	63
	.byte	25
	.byte	0
	.byte	0
	.byte	3
	.byte	5
	.byte	0
	.byte	2
	.byte	24
	.byte	3
	.byte	37
	.byte	58
	.byte	11
	.byte	59
	.byte	11
	.byte	73
	.byte	19
	.byte	0
	.byte	0
	.byte	4
	.byte	46
	.byte	0
	.byte	17
	.byte	27
	.byte	18
	.byte	6
	.byte	64
	.byte	24
	.byte	3
	.byte	37
	.byte	58
	.byte	11
	.byte	59
	.byte	11
	.byte	63
	.byte	25
	.byte	0
	.byte	0
	.byte	5
	.byte	46
	.byte	1
	.byte	17
	.byte	27
	.byte	18
	.byte	6
	.byte	64
	.byte	24
	.byte	3
	.byte	37
	.byte	58
	.byte	11
	.byte	59
	.byte	11
	.byte	73
	.byte	19
	.byte	63
	.byte	25
	.byte	0
	.byte	0
	.byte	6
	.byte	52
	.byte	0
	.byte	2
	.byte	24
	.byte	3
	.byte	37
	.byte	58
	.byte	11
	.byte	59
	.byte	11
	.byte	73
	.byte	19
	.byte	0
	.byte	0
	.byte	7
	.byte	36
	.byte	0
	.byte	3
	.byte	37
	.byte	62
	.byte	11
	.byte	11
	.byte	11
	.byte	0
	.byte	0
	.byte	0
	.section	.debug_info,"",@progbits
.Lcu_begin0:
	.word	.Ldebug_info_end0-.Ldebug_info_start0
.Ldebug_info_start0:
	.half	5
	.byte	1
	.byte	4
	.word	.debug_abbrev
	.byte	1
	.byte	0
	.half	29
	.byte	1
	.word	.Lstr_offsets_base0
	.word	.Lline_table_start0
	.byte	2
	.byte	0
	.word	.Lfunc_end2-.Lfunc_begin0
	.word	.Laddr_table_base0
	.byte	2
	.byte	0
	.word	.Lfunc_end0-.Lfunc_begin0
	.byte	1
	.byte	88
	.byte	3
	.byte	0
	.byte	4

	.word	133

	.byte	3
	.byte	2
	.byte	145
	.byte	116
	.byte	8
	.byte	0
	.byte	4
	.word	137
	.byte	0
	.byte	4
	.byte	1
	.word	.Lfunc_end1-.Lfunc_begin1
	.byte	1
	.byte	88
	.byte	5
	.byte	0
	.byte	7

	.byte	5
	.byte	2
	.word	.Lfunc_end2-.Lfunc_begin2
	.byte	1
	.byte	88
	.byte	6
	.byte	0
	.byte	15
	.word	137

	.byte	6
	.byte	2
	.byte	145
	.byte	116
	.byte	8
	.byte	0
	.byte	16
	.word	137
	.byte	6
	.byte	2
	.byte	145
	.byte	112
	.byte	9
	.byte	0
	.byte	16
	.word	137
	.byte	6
	.byte	2
	.byte	145
	.byte	111
	.byte	10
	.byte	0
	.byte	17
	.word	133
	.byte	6
	.byte	2
	.byte	145
	.byte	110
	.byte	11
	.byte	0
	.byte	17
	.word	133
	.byte	0
	.byte	7
	.byte	4
	.byte	8
	.byte	1
	.byte	7
	.byte	7
	.byte	5
	.byte	4
	.byte	0
.Ldebug_info_end0:
	.section	.debug_str_offsets,"",@progbits
	.word	52
	.half	5
	.half	0
.Lstr_offsets_base0:
	.section	.debug_str,"MS",@progbits,1
.Linfo_string0:
	.asciz	"Homebrew clang version 16.0.6"
.Linfo_string1:
	.asciz	"main.c"
.Linfo_string2:
	.asciz	"/Users/scott/repo/sokoide/rv32i-go/data"
.Linfo_string3:
	.asciz	"is_even"
.Linfo_string4:
	.asciz	"char"
.Linfo_string5:
	.asciz	"riscv32_boot"
.Linfo_string6:
	.asciz	"main"
.Linfo_string7:
	.asciz	"int"
.Linfo_string8:
	.asciz	"a"
.Linfo_string9:
	.asciz	"b"
.Linfo_string10:
	.asciz	"c"
.Linfo_string11:
	.asciz	"d"
	.section	.debug_str_offsets,"",@progbits
	.word	.Linfo_string0
	.word	.Linfo_string1
	.word	.Linfo_string2
	.word	.Linfo_string3
	.word	.Linfo_string4
	.word	.Linfo_string5
	.word	.Linfo_string6
	.word	.Linfo_string7
	.word	.Linfo_string8
	.word	.Linfo_string9
	.word	.Linfo_string10
	.word	.Linfo_string11
	.section	.debug_addr,"",@progbits
	.word	.Ldebug_addr_end0-.Ldebug_addr_start0
.Ldebug_addr_start0:
	.half	5
	.byte	4
	.byte	0
.Laddr_table_base0:
	.word	.Lfunc_begin0
	.word	.Lfunc_begin1
	.word	.Lfunc_begin2
.Ldebug_addr_end0:
	.ident	"Homebrew clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	.addrsig_sym is_even
	.addrsig_sym main
	.addrsig_sym _out
	.section	.debug_line,"",@progbits
.Lline_table_start0:
