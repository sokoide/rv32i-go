package main

import (
	"strconv"
)

func evaluate_stmt(stmt statement) error {
	return nil
}

func evaluate_expr(expr expression) (int, error) {
	switch e := expr.(type) {

	case *numberExpression:
		v, err := strconv.Atoi(e.Lit)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *parenExpression:
		v, err := evaluate_expr(e.SubExpr)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *binOpExpression:
		lhsV, err := evaluate_expr(e.LHS)
		if err != nil {
			return 0, err
		}
		rhsV, err := evaluate_expr(e.RHS)
		if err != nil {
			return 0, err
		}
		switch e.Operator {
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
