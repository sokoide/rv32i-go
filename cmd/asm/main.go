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
	addi	sp, sp, -12
	add	sp, sp, a0
	jal riscv32_boot
_out:
	ret
is_even:
	addi	sp, sp, -16
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	addi	s0, sp, 16
	sw	a0, -12(s0)
	lw	a0, -12(s0)
	srli	a1, a0, 31
	add	a1, a0, a1
	andi	a1, a1, -2
	sub	a0, a0, a1
	seqz a0, a0
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	addi	sp, sp, 16
	ret
riscv32_boot:
	addi	sp, sp, -16
	sw	ra, 12(sp)
	sw	s0, 8(sp)
	addi	s0, sp, 16
	auipc	ra, 0
	jalr	24(ra)
	lw	ra, 12(sp)
	lw	s0, 8(sp)
	addi	sp, sp, 16
	ret
`

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
