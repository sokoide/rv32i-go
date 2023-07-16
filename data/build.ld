ENTRY(boot)

SECTIONS {
    /* program origin */
    . = 0x00000000;
    __ram_start = .;

    .text : {

        __text = .;
        KEEP(*(.boot));

        . = ALIGN(16);
        *(.text .text.*);

        . = ALIGN(4096);
        __text_end = .;
    }

    __data = .;
    .rodata : {
        *(.rodata .rodata.*);
        . = ALIGN(16);
        *(.srodata .srodata.*);

        KEEP(*(.symbols));

        . = ALIGN(4096);
        __boot_elf = .;
        *(.boot_elf);
        . = ALIGN(4096);
    }

    .data : {
        . = ALIGN(16);
        *(.data .data.*);
        . = ALIGN(16);
        *(.sdata .sdata.*);
    }

    .bss : {
        __bss = .;
        . = ALIGN(16);
        *(.bss .bss.*);
        . = ALIGN(16);
        *(.sbss .sbss.*);
        __bss_end = .;

        . = ALIGN(4096);

    /* stack 16KB */
        __boot_stack = .;
        . = . + 4 * (16 * 1024);
    }

    . = ALIGN(4096);
    __data_end = .;
    __free_ram_start = .;
}