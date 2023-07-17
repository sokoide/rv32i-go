.align 4
.global boot
.global _out

// Binary entry defined in build.ld
boot:
	// clear ra and fp to stop stacktrace here
	mv ra, zero
	mv fp, zero

	// stack defined in build.ld
	// initial sp: __boot_stack + 16*1024
	// stack upper limit: __boot_stack
	li a0, 16*1024
	la sp, __boot_stack
	// sp = __boot_stack + a0
	add sp, sp, a0
	jal riscv32_boot

_out:
	ret
