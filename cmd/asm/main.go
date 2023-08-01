package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
	"github.com/sokoide/rv32i-go/pkg/rv32iasm"
)

func main() {
	var err error

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
_out:
	ret`

	log.Tracef("src: %s", src)

	reader := strings.NewReader(src)
	scanner := rv32iasm.NewScanner(reader)

	var program *rv32iasm.Program
	program, err = scanner.Parse()
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	log.Debugf("* program=%+v", program)

	log.Info("* start evaluation")
	ev := rv32iasm.NewEvaluator()
	err = ev.EvaluateProgram(program)
	if err != nil {
		log.Fatalf("EvaluateProgram error: %v", err)
	}

	emu := rv32i.NewEmulator()
	emu.Reset()
}
