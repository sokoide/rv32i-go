.align 4
.global boot
.global _out

boot:
    // clear ra and fp to stop stacktrace here
    mv ra, zero
    mv fp, zero

        // initial sp: __boot_stack + 16*1024
        // stack upper limit: __boot_stack
        li a0, 16*1024
    la sp, __boot_stack
        add sp, sp, a0            // sp = __boot_stack + a0

    jal riscv32_boot

_out:
        ret