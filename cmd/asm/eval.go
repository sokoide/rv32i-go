package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

type Evaluator struct {
	labels map[string]int
	code   []uint32
	PC     int
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) reset() {
	e.labels = make(map[string]int, 0)
	e.code = make([]uint32, 0)
	e.PC = 0
}

func (e *Evaluator) evaluate_program(prog *program) error {
	log.Infof("evaluate_program: stmt=%d", len(prog.statements))
	e.reset()

	for idx, stmt := range prog.statements {
		codes, generated := e.gen_code(stmt)
		if generated {
			for _, code := range codes {
				log.Debugf("[%d] %+v, 0x%08x, 0b%032b\n", idx, stmt, code, code)
				e.code = append(e.code, code)
			}
			e.PC += 4 * len(codes)
		}
	}

	// dumps")
	log.Debug("Labels)")
	for key, val := range e.labels {
		log.Debugf("%16s: 0x%08x", key, val)
	}
	log.Debug("Code)")
	for idx, code := range e.code {
		log.Debugf("0x%08x: 0x%08x", idx*4, code)
	}
	return nil
}

func (e *Evaluator) gen_code(stmt *statement) ([]uint32, bool) {
	switch stmt.opcode {
	case "lui":
		// op1: rd, op2, imm
		return []uint32{rv32i.GenCode(rv32i.OpLui, stmt.op1, stmt.op2, stmt.op3)}, true
	case "auipc":
		// op1: rd, op2: imm
		return []uint32{rv32i.GenCode(rv32i.OpAuipc, stmt.op1, stmt.op2, stmt.op3)}, true
	case "addi":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpAddi, stmt.op1, stmt.op2, stmt.op3)}, true
	case "li":
		// op1: rd, op2: imm
		if (stmt.op2 & 0b01111111_11111111_11111000_00000000) != 0 {
			return []uint32{rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, stmt.op2)}, true
		} else {
			hi := int(rv32i.SignExtension((uint32(stmt.op2) >> 12), 20))
			low := stmt.op2 & 0b1111_1111_1111
			log.Debugf("%d, %d\n", hi, low)
			return []uint32{
				rv32i.GenCode(rv32i.OpLui, stmt.op1, hi, 0),
				rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, low),
			}, true
		}
	case "add":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpAdd, stmt.op1, stmt.op2, stmt.op3)}, true
	case "jal":
		if stmt.str1 == "" {
			// op1: rd, op2: offset
			return []uint32{rv32i.GenCode(rv32i.OpJal, stmt.op1, stmt.op2, stmt.op3)}, true
		} else {
			// op1: label
			// if the labels is located  +/- 1KB from regx, use 'jalr regx, imm'
			// or, if the absolute address is <512KB from regx, use 'jal regx, imm'
			// otherwise, insert auipc and jal
			// to simplify, the assembler always uses 'jal x0, imm' assuming the target is
			// located between 0x0000 and 0x80000 (512KB)
			panic("jal label not implemented")
		}
	case "label":
		e.labels[stmt.str1] = e.PC
		return []uint32{0}, false
	case "comment":
		return []uint32{0}, false
	default:
		// TODO:
		return []uint32{0}, false
	}
}

func (e *Evaluator) evaluate_expr(expr expression) (int, error) {
	switch ex := expr.(type) {

	case *numberExpression:
		v, err := strconv.Atoi(ex.Lit)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *parenExpression:
		v, err := e.evaluate_expr(ex.SubExpr)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *binOpExpression:
		lhsV, err := e.evaluate_expr(ex.LHS)
		if err != nil {
			return 0, err
		}
		rhsV, err := e.evaluate_expr(ex.RHS)
		if err != nil {
			return 0, err
		}
		switch ex.Operator {
		case '+':
			return lhsV + rhsV, nil
		case '-':
			return lhsV - rhsV, nil
		case '*':
			return lhsV * rhsV, nil
		case '/':
			return lhsV / rhsV, nil
		case '%':
			return lhsV % rhsV, nil
		default:
			panic("Unknown operator")
		}
	default:
		panic("Unknown Expression type")
	}
}
