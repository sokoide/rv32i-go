%{
package rv32iasm

import (
    "strconv"
    log "github.com/sirupsen/logrus"
)
func chkerr(err error) {
    if err != nil {
        panic(err)
    }
}
%}

%union{
    program *Program
    stmt    *statement
    expr    expression
    tok     token
}

%type<program> program
%type<stmt> stmt lui_stmt auipc_stmt addi_stmt li_stmt
%type<stmt> add_stmt jal_stmt
%type<stmt> label_stmt
%type<expr> expr
%token<tok> LF COLON COMMA NUMBER IDENT
%token<tok> REGISTER AUIPC LUI ADDI ADD LI JAL

%left '+' '-'
%left '*' '/'

%start program

%%
program: /* empty */ {
    log.Debug("* empty program")
        $$ = &Program{
            statements: make([]*statement, 0),
        }
        assemblerlex.(*lexer).program = $$
    }
    | program stmt LF {
        log.Debugf("* appendind stmt %v, stmt count %d", $2, len($$.statements))
        $$ = &Program {
            statements: append($1.statements, $2),
        }
        assemblerlex.(*lexer).program = $$
    }

stmt: /* empty */ {
        log.Debug("* comment or empty stmt")
        $$ = &statement{
            opcode: "comment",
        }
    }
    | lui_stmt { $$ = $1 }
    | auipc_stmt { $$ = $1 }
    | addi_stmt { $$ = $1 }
    | li_stmt { $$ = $1 }
    | add_stmt { $$ = $1 }
    | jal_stmt { $$ = $1 }
    | label_stmt { $$ = $1}
    | expr {
        log.Debugf("* stmt expr %v", $$)
        $$ = &statement{
            opcode: "expr",
        }
    }

lui_stmt: LUI REGISTER COMMA NUMBER {
        log.Debugf("* lui_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
        }
    }

auipc_stmt: AUIPC REGISTER COMMA NUMBER {
        log.Debugf("* auipc_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
        }
    }

addi_stmt: ADDI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* addi_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: val,
        }
    }

li_stmt: LI REGISTER COMMA NUMBER {
        log.Debugf("* li_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
        }
    }

add_stmt: ADD REGISTER COMMA REGISTER COMMA REGISTER {
        log.Debugf("* add_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: regs[$6.lit],
        }
    }

jal_stmt: JAL REGISTER COMMA NUMBER {
        log.Debugf("* jal_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
        }
    }
    | JAL IDENT {
        log.Debugf("* jal_stmt (label): %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            str1: $2.lit,
        }
    }


label_stmt: IDENT COLON {
        log.Debugf("* label_stmt: %+v", $1)
        $$ = &statement{
            opcode: "label",
            str1: $1.lit,
        }
    }

expr: NUMBER {
        $$ = &numberExpression{Lit: $1.lit}
	}
    | expr '+' expr {
        $$ = &binOpExpression{LHS: $1, Operator: int('+'), RHS: $3}
	}
    | expr '-' expr {
        $$ = &binOpExpression{LHS: $1, Operator: int('-'), RHS: $3}
	}
    | expr '*' expr {
        $$ = &binOpExpression{LHS: $1, Operator: int('*'), RHS: $3}
	}
    | expr '/' expr {
        $$ = &binOpExpression{LHS: $1, Operator: int('/'), RHS: $3}
	}
    | '(' expr ')' {
        $$ = &parenExpression{SubExpr: $2}
    }

%%
