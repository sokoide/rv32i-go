package rv32iasm

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

type resolveTarget struct {
	op     string
	symbol string
}

type Evaluator struct {
	labels         map[string]int
	linksToResolve map[int]resolveTarget
	Code           []uint32
	PC             int
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Reset() {
	e.labels = make(map[string]int, 0)
	e.linksToResolve = make(map[int]resolveTarget, 0)
	e.Code = make([]uint32, 0)
	e.PC = 0
}

func (e *Evaluator) Assemble(reader io.Reader) ([]string, error) {
	var err error
	var program *Program
	scanner := NewScanner(reader)

	program, err = scanner.Parse()
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	log.Debugf("* program=%+v", program)

	log.Debug("* start evaluation")
	ev := NewEvaluator()

	return ev.EvaluateProgram(program)
}

func (e *Evaluator) EvaluateProgram(prog *Program) ([]string, error) {
	s := make([]string, 0)
	log.Debugf("EvaluateProgram: stmt=%d", len(prog.statements))
	e.Reset()

	for idx, stmt := range prog.statements {
		codes, generated := e.gen_code(stmt)
		if generated {
			for _, code := range codes {
				log.Debugf("[%d] %+v, 0x%08x, 0b%032b\n", idx, stmt, code, code)
				e.Code = append(e.Code, code)
			}
			e.PC += 4 * len(codes)
		}
	}

	// dump
	labelsR := make(map[uint32]string, 0)

	log.Debug("Labels)")
	for key, val := range e.labels {
		log.Debugf("%-16s: 0x%08x", key, val)
		labelsR[uint32(val)] = key
	}

	log.Debug("Links to Resolve)")
	for key, val := range e.linksToResolve {
		log.Debugf("0x%08x: %s, %s", key, val.op, val.symbol)
	}

	e.resolveLinks()
	log.Debug("After resolved)")
	for idx, code := range e.Code {
		instr := rv32i.NewInstruction(code)
		addr := uint32(idx * 4)
		if val, ok := labelsR[addr]; ok {
			log.Debugf("%08x <%s>:", addr, val)
			s = append(s, fmt.Sprintf("%08x <%s>:", addr, val))
		}
		log.Debugf("%8x: 0x%08x %s", addr, code, instr.GetCodeString())
		s = append(s, fmt.Sprintf("%8x: 0x%08x %s", addr, code, instr.GetCodeString()))
	}

	return s, nil
}

func (e *Evaluator) gen_code(stmt *statement) ([]uint32, bool) {
	switch stmt.opcode {
	case "lui":
		// op1: rd, op2, imm
		return []uint32{rv32i.GenCode(rv32i.OpLui, stmt.op1, stmt.op2, stmt.op3)}, true
	case "auipc":
		// op1: rd, op2: imm
		return []uint32{rv32i.GenCode(rv32i.OpAuipc, stmt.op1, stmt.op2, stmt.op3)}, true
	case "jal":
		if stmt.str1 == "" {
			// op1: rd, op2: offset
			return []uint32{rv32i.GenCode(rv32i.OpJal, stmt.op1, stmt.op2, stmt.op3)}, true
		} else {
			// str1: label
			if val, ok := e.labels[stmt.str1]; ok {
				imm := val - e.PC
				return []uint32{rv32i.GenCode(rv32i.OpJal, stmt.op1, imm, 0)}, true
			} else {
				e.linksToResolve[e.PC] = resolveTarget{"jal", stmt.str1}
				return []uint32{rv32i.GenCode(rv32i.OpJal, stmt.op1, 0, 0)}, true
			}
		}
	case "jalr":
		// op1: rd, op2: offset, op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpJalr, stmt.op1, stmt.op2, stmt.op3)}, true
	case "beq":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBeq, stmt.op1, stmt.op2, stmt.op3)}, true
	case "bne":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBne, stmt.op1, stmt.op2, stmt.op3)}, true
	case "blt":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBlt, stmt.op1, stmt.op2, stmt.op3)}, true
	case "bge":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBge, stmt.op1, stmt.op2, stmt.op3)}, true
	case "bltu":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBltu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "bgeu":
		// op1: rs1, op2: rs2, op3: offset
		return []uint32{rv32i.GenCode(rv32i.OpBgeu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "lb":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpLb, stmt.op1, stmt.op2, stmt.op3)}, true
	case "lh":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpLh, stmt.op1, stmt.op2, stmt.op3)}, true
	case "lw":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpLw, stmt.op1, stmt.op2, stmt.op3)}, true
	case "lbu":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpLbu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "lhu":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpLhu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sb":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpSb, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sh":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpSh, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sw":
		// op1: rs2, op2: offset: op3: rs1
		return []uint32{rv32i.GenCode(rv32i.OpSw, stmt.op1, stmt.op2, stmt.op3)}, true
	case "addi":
		// op1: rd, op2: rs1, op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpAddi, stmt.op1, stmt.op2, stmt.op3)}, true
	case "slti":
		// op1: rd, op2: rs1, op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpSlti, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sltiu":
		// op1: rd, op2: rs1, op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpSltiu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "xori":
		// op1: rd, op2: rs1, op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpXori, stmt.op1, stmt.op2, stmt.op3)}, true
	case "ori":
		// op1: rd, op2: rs1, op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpOri, stmt.op1, stmt.op2, stmt.op3)}, true
	case "andi":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpAndi, stmt.op1, stmt.op2, stmt.op3)}, true
	case "slli":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpSlli, stmt.op1, stmt.op2, stmt.op3)}, true
	case "srli":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpSrli, stmt.op1, stmt.op2, stmt.op3)}, true
	case "srai":
		// op1: rd, op2: rs1: op3: imm
		return []uint32{rv32i.GenCode(rv32i.OpSrai, stmt.op1, stmt.op2, stmt.op3)}, true
	case "add":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpAdd, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sub":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSub, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sll":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSll, stmt.op1, stmt.op2, stmt.op3)}, true
	case "slt":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSlt, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sltu":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSltu, stmt.op1, stmt.op2, stmt.op3)}, true
	case "xor":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpXor, stmt.op1, stmt.op2, stmt.op3)}, true
	case "srl":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSrl, stmt.op1, stmt.op2, stmt.op3)}, true
	case "sra":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpSra, stmt.op1, stmt.op2, stmt.op3)}, true
	case "or":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpOr, stmt.op1, stmt.op2, stmt.op3)}, true
	case "and":
		// op1: rd, op2: rs1: op3: rs2
		return []uint32{rv32i.GenCode(rv32i.OpAnd, stmt.op1, stmt.op2, stmt.op3)}, true
	// pseudo instructions
	case "call":
		// op1: rd, str1: symbol
		if val, ok := e.labels[stmt.str1]; ok {
			hi := int(rv32i.SignExtension((uint32(val) >> 12), 20))
			low := val & 0b1111_1111_1111
			return []uint32{
				rv32i.GenCode(rv32i.OpAuipc, stmt.op1, hi, 0),
				rv32i.GenCode(rv32i.OpJalr, stmt.op1, low, 0),
			}, true
		} else {
			e.linksToResolve[e.PC] = resolveTarget{"call", stmt.str1}
			return []uint32{
				rv32i.GenCode(rv32i.OpAuipc, stmt.op1, 0, 0),
				rv32i.GenCode(rv32i.OpJalr, stmt.op1, 0, 0),
			}, true
		}
	case "la":
		// op1: rd, str1: symbol
		if val, ok := e.labels[stmt.str1]; ok {
			hi := int(rv32i.SignExtension((uint32(val) >> 12), 20))
			low := val & 0b1111_1111_1111
			return []uint32{
				rv32i.GenCode(rv32i.OpAuipc, stmt.op1, hi, 0),
				rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, low),
			}, true
		} else {
			e.linksToResolve[e.PC] = resolveTarget{"la", stmt.str1}
			return []uint32{
				rv32i.GenCode(rv32i.OpAuipc, stmt.op1, 0, 0),
				rv32i.GenCode(rv32i.OpAddi, stmt.op1, 0, 0),
			}, true
		}
	case "li":
		// op1: rd, op2: imm
		if (stmt.op2 & 0b01111111_11111111_11111000_00000000) == 0 {
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

func (e *Evaluator) resolveLinks() error {
	for PC, rt := range e.linksToResolve {
		if val, ok := e.labels[rt.symbol]; ok {
			switch rt.op {
			case "call":
				rd := int((e.Code[PC/4] >> 7) & 0b11111)
				hi := int(rv32i.SignExtension((uint32(val) >> 12), 20))
				low := val & 0b1111_1111_1111
				e.Code[PC/4] = rv32i.GenCode(rv32i.OpAuipc, rd, hi, 0)
				e.Code[PC/4+1] = rv32i.GenCode(rv32i.OpJalr, rd, low, 0)
			case "jal":
				// if the label is located +/-512KB from regx, use 'jal regx, imm'
				// TODO: otherwise, insert auipc and jal
				// To simplify, the assembler always uses 'jal x1, imm' assuming the target is
				// located between 0x0000 and 0x80000 (512KB)
				if rv32i.Abs(val-PC) <= 512*1024 {
					imm := val - PC
					rd := int((e.Code[PC/4] >> 7) & 0b11111)
					e.Code[PC/4] = rv32i.GenCode(rv32i.OpJal, rd, imm, 0)
				} else {
					return fmt.Errorf("label %s is too far! PC:%x, %s:%x", rt.symbol, PC, rt.symbol, val)
				}
			case "la":
				rd := int((e.Code[PC/4] >> 7) & 0b11111)
				hi := int(rv32i.SignExtension((uint32(val) >> 12), 20))
				low := val & 0b1111_1111_1111
				e.Code[PC/4] = rv32i.GenCode(rv32i.OpAuipc, rd, hi, 0)
				e.Code[PC/4+1] = rv32i.GenCode(rv32i.OpAddi, rd, 0, low)
			}
		} else {
			return fmt.Errorf("label %s not found", rt.symbol)
		}
	}
	return nil
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
			return 0, errors.New("Unknown operator")
		}
	default:
		return 0, errors.New("Unknown Expression type")
	}
}
