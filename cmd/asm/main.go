package main

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
	"github.com/sokoide/rv32i-go/pkg/rv32iasm"
)

func main() {
	ev := rv32iasm.NewEvaluator()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	// log.SetLevel(log.InfoLevel)

	log.Info("asm started")
	src := `boot:
# This is a comment line
	li ra, 0
	li s0, 0 # This is a comment
	lui a0, 4
	auipc sp, 1
	addi sp, sp, -12
	add sp, sp, a0
	jal riscv32_boot
	li ra, -300 # This is never called
	li a1, 1000000000 # This is never called
	li a0, 1 # This is never called
	li a1, 2 # This is never called
	li a3, 3 # This is never called
riscv32_boot:
	addi	sp, sp, -16
	li ra, 1
	jal boot
`

	log.Tracef("src: %s", src)

	s := bufio.NewScanner(strings.NewReader(src))
	scanner := new(rv32iasm.Scanner)
	source := []string{}
	for s.Scan() {
		source = append(source, s.Text())
	}
	scanner.Init(strings.Join(source, "\n") + "\n")

	var program *rv32iasm.Program = scanner.Parse()
	log.Debugf("* program=%+v", program)

	log.Info("* start evaluation")
	err := ev.EvaluateProgram(program)
	if err != nil {
		panic(nil)
	}

	emu := rv32i.NewEmulator()
	emu.Reset()
}
