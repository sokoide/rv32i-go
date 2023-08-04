# RV32I emulator

## How to build

```sh
make
```

## How to test

```sh
make test
```

## How to run the assembler

```sh
make
./asm -source data/demo.s
```

## How to run the assembled code with the emulator

* Please take a look at [cmd/demo/main.go](./cmd/demo/main.go)

## RV32I registers

```txt
// Register ABIName Description                         Saver
// -----------------------------------------------------------
// x0       zero    Hard-wired zero                     —
// x1       ra      Return address                      Caller
// x2       sp      Stack pointer                       Callee
// x3       gp      Global pointer                      —
// x4       tp      Thread pointer                      —
// x5–7     t0–2    Temporaries                         Caller
// x8       s0/fp   Saved register/frame pointer        Callee
// x9       s1      Saved register                      Callee
// x10–11   a0–1    Function arguments/return values    Caller
// x12–17   a2–7    Function arguments                  Caller
// x18–27   s2–11   Saved registers                     Callee
// x28–31   t3–6    Temporaries                         Caller
```

## Supported Instructions

### Regular Instructions

* Major RV32I instructions are supported except for fence*, ecall, ebreak and cs*

### Pesudo Instructions

* `li`, `call`, `ret` or some other limited pseudo instructions are supported
* Please refer to [assembler.y](./pkg/rv32iasm/assembler.y) for the complete list of the supported instructions

# RV32I instructions

* [RV32I, RV64I instructions](https://msyksphinz-self.github.io/riscv-isadoc/html/rvi.html)

## How to make RV32I assembly for local testing

* `main.c`, `build.ld` and `start.S` are built into RV32I by `Makefile` in `data` directory
* Install llvm by `brew install llvm`. Apple's default LLVM doens't generate code for riscv32 target

```sh
$ make
mkdir -p ~/tmp/riscv1
clang -std=c11 -Wall -g3 -O0 --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -ffreestanding -fno-builtin -c start.S -o /Users/scott/tmp/riscv1/start.o
clang -std=c11 -Wall -g3 -O0 --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -ffreestanding -fno-builtin -c main.c -o /Users/scott/tmp/riscv1/main.o
clang /Users/scott/tmp/riscv1/start.o /Users/scott/tmp/riscv1/main.o -o ~/tmp/riscv1/riscv1 -static --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -Tbuild.ldn

$ make disass
clang /Users/scott/tmp/riscv1/start.o /Users/scott/tmp/riscv1/main.o -o ~/tmp/riscv1/riscv1 -static --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -Tbuild.ld
llvm-objdump -D ~/tmp/riscv1/riscv1

/Users/scott/tmp/riscv1/riscv1: file format elf32-littleriscv

Disassembly of section .text:

00000000 <boot>:
       0: 93 00 00 00   li      ra, 0
       4: 13 04 00 00   li      s0, 0
       8: 37 45 00 00   lui     a0, 4

0000000c <.Lpcrel_hi0>:
       c: 17 11 00 00   auipc   sp, 1
      10: 13 01 41 ff   addi    sp, sp, -12
      14: 33 01 a1 00   add     sp, sp, a0
      18: ef 00 40 04   jal     0x5c <riscv32_boot>

0000001c <_out>:
      1c: 67 80 00 00   ret

00000020 <is_even>:
      20: 13 01 01 ff   addi    sp, sp, -16
      24: 23 26 11 00   sw      ra, 12(sp)
      28: 23 24 81 00   sw      s0, 8(sp)
      2c: 13 04 01 01   addi    s0, sp, 16
      30: 23 2a a4 fe   sw      a0, -12(s0)
      34: 03 25 44 ff   lw      a0, -12(s0)
      38: 93 55 f5 01   srli    a1, a0, 31
      3c: b3 05 b5 00   add     a1, a0, a1
      40: 93 f5 e5 ff   andi    a1, a1, -2
      44: 33 05 b5 40   sub     a0, a0, a1
      48: 13 35 15 00   seqz    a0, a0
      4c: 83 20 c1 00   lw      ra, 12(sp)
...
```

## Execution Example

* `data` directory has some examples
* `make test` runs `sample-binary-*.txt` files
* `make run` runs an embedded sample 3 instruction binary and dump registers, then do the same for `sample-binary-003.bin`

```sh
make run
INFO[0000] * Started
INFO[0000] * Running a small program

# this runs the following 3 instructions
    // test binary in sample-binary-001.txt
    // 80000000 <boot>:
    // 80000000: 93 00 00 00   li      ra, 0
    // # x1 == ra == 0
    // 80000004: 13 04 00 00   li      s0, 0
    // # x8 == s0/fp == 0
    // 80000008: 37 45 00 00   lui     a0, 4
    // # x10 == a0 == 4<<12 == 16384

# Then dump registers
INFO[0000] * Registers
INFO[0000] x0 = 0, 0x00000000
INFO[0000] x1 = 0, 0x00000000
INFO[0000] x2 = 0, 0x00000000
INFO[0000] x3 = 0, 0x00000000
INFO[0000] x4 = 0, 0x00000000
INFO[0000] x5 = 0, 0x00000000
INFO[0000] x6 = 0, 0x00000000
INFO[0000] x7 = 0, 0x00000000
INFO[0000] x8 = 0, 0x00000000
INFO[0000] x9 = 0, 0x00000000
INFO[0000] x10 = 16384, 0x00004000
INFO[0000] x11 = 0, 0x00000000
INFO[0000] x12 = 0, 0x00000000
INFO[0000] x13 = 0, 0x00000000
INFO[0000] x14 = 0, 0x00000000
INFO[0000] x15 = 0, 0x00000000
INFO[0000] x16 = 0, 0x00000000
INFO[0000] x17 = 0, 0x00000000
INFO[0000] x18 = 0, 0x00000000
INFO[0000] x19 = 0, 0x00000000
INFO[0000] x20 = 0, 0x00000000
INFO[0000] x21 = 0, 0x00000000
INFO[0000] x22 = 0, 0x00000000
INFO[0000] x23 = 0, 0x00000000
INFO[0000] x24 = 0, 0x00000000
INFO[0000] x25 = 0, 0x00000000
INFO[0000] x26 = 0, 0x00000000
INFO[0000] x27 = 0, 0x00000000
INFO[0000] x28 = 0, 0x00000000
INFO[0000] x29 = 0, 0x00000000
INFO[0000] x30 = 0, 0x00000000
INFO[0000] x31 = 0, 0x00000000
INFO[0000] pc = 0x0000000c

# This converts sample-binary-003.txt into sample-binary-003.bin
INFO[0000] * Converting txt to bin

INFO[0000] * Running the bin

# Then dump registers
INFO[0000] * Registers
INFO[0000] x0 = 92, 0x0000005c
INFO[0000] x1 = 208, 0x000000d0
INFO[0000] x2 = 20432, 0x00004fd0
INFO[0000] x3 = 0, 0x00000000
INFO[0000] x4 = 0, 0x00000000
INFO[0000] x5 = 0, 0x00000000
INFO[0000] x6 = 0, 0x00000000
INFO[0000] x7 = 0, 0x00000000
INFO[0000] x8 = 20464, 0x00004ff0
INFO[0000] x9 = 0, 0x00000000
INFO[0000] x10 = 10, 0x0000000a
INFO[0000] x11 = 0, 0x00000000
INFO[0000] x12 = 0, 0x00000000
INFO[0000] x13 = 0, 0x00000000
INFO[0000] x14 = 0, 0x00000000
INFO[0000] x15 = 0, 0x00000000
INFO[0000] x16 = 0, 0x00000000
INFO[0000] x17 = 0, 0x00000000
INFO[0000] x18 = 0, 0x00000000
INFO[0000] x19 = 0, 0x00000000
INFO[0000] x20 = 0, 0x00000000
INFO[0000] x21 = 0, 0x00000000
INFO[0000] x22 = 0, 0x00000000
INFO[0000] x23 = 0, 0x00000000
INFO[0000] x24 = 0, 0x00000000
INFO[0000] x25 = 0, 0x00000000
INFO[0000] x26 = 0, 0x00000000
INFO[0000] x27 = 0, 0x00000000
INFO[0000] x28 = 0, 0x00000000
INFO[0000] x29 = 0, 0x00000000
INFO[0000] x30 = 0, 0x00000000
INFO[0000] x31 = 0, 0x00000000
INFO[0000] pc = 0x0000001c
INFO[0000] * Completed
```
