package main

type (
	program struct {
		statements []*statement
	}

	statement struct {
		opcode string
		op1    int
		op2    int
		op3    int
	}
	expression interface {
		expression()
	}
)

type (
	numberExpression struct {
		Lit string
	}

	parenExpression struct {
		SubExpr expression
	}

	binOpExpression struct {
		LHS      expression
		Operator int
		RHS      expression
	}
)

func (x *numberExpression) expression() {}
func (x *parenExpression) expression()  {}
func (x *binOpExpression) expression()  {}
