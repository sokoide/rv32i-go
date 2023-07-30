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
	case "addi":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpAddi, stmt.op1, stmt.op2, stmt.op3)}, true
	case "li":
		// op1: rd, op2: imm
		if (stmt.op2 & 0b01111111_11111111_11111000_00000000) != 0 {
			return []uint32{rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, stmt.op2)}, true
		} else {
			// TODO: support negative numbers
			log.Debugf("%d, %d\n", stmt.op2>>12, stmt.op2&0b1111_1111_1111)
			return []uint32{
				rv32i.GenCode(rv32i.OpLui, stmt.op1, stmt.op2>>12, 0),
				rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, stmt.op2&0b1111_1111_1111),
			}, true
		}
	case "lui":
		// op1: rd, op2, imm
		return []uint32{rv32i.GenCode(rv32i.OpLui, stmt.op1, stmt.op2, stmt.op3)}, true
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
