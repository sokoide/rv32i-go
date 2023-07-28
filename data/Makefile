#
# Makefile
# install llvm by brew install llvm and use the clang in it
#
# 'make disass' to generate the sample-binary-*
#

TARGET = riscv1
OUTDIR = ~/tmp/riscv1

COMMONFILES =
SRCS = main.c
ASMSRCS = start.s

HEADERS :=  $(shell find . -type f -name '*.h')
OBJS := $(ASMSRCS:.s=.o) $(SRCS:.c=.o)

OUTOBJS = $(addprefix $(OUTDIR)/,$(OBJS))

CC = clang
CCFLAGS = -std=c11 -Wall -g3 -O0 --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -ffreestanding -fno-builtin
LDFLAGS = -static --target=riscv32 -march=rv32i -mabi=ilp32 -mno-relax -nostdlib -Tbuild.ld

default: $(TARGET)

$(TARGET): build

build: $(OUTOBJS)
	$(CC) $^ -o $(OUTDIR)/$(TARGET) $(LDFLAGS)

$(OUTDIR)/%.o : %.c | $(OUTDIR)
	$(CC) $(CCFLAGS) -c $< -o $@

$(OUTDIR)/%.o : %.s | $(OUTDIR)
	$(CC) $(CCFLAGS) -c $< -o $@

$(OUTDIR):
	mkdir -p $(OUTDIR)

debug: $(TARGET)
	lldb $(OUTDIR)/$(TARGET)

diag: $(TARGET)
	llvm-readelf -a $(OUTDIR)/$(TARGET)
	llvm-objdump -p $(OUTDIR)/$(TARGET)

disass: $(TARGET)
	llvm-objdump -D $(OUTDIR)/$(TARGET)

disassstart: $(TARGET)
	llvm-objdump -D $(OUTDIR)/start.o

clean:
	rm -rf $(OUTDIR)
