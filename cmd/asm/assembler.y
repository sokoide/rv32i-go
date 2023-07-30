%{
package main

import (
    "fmt"
    "strconv"
)

%}

%union{
    program *program
    stmt    *statement
    expr    expression
    tok     token
}

%type<program> program
%type<stmt> stmt lui_stmt auipc_stmt addi_stmt li_stmt
%type<stmt> add_stmt label_stmt
%type<expr> expr
%token<tok> LF COLON COMMA NUMBER IDENT
%token<tok> REGISTER AUIPC LUI ADDI ADD LI

%left '+' '-'
%left '*' '/'

%start program

%%
program: /* empty */ {
    fmt.Println("* empty program")
        $$ = &program{
            statements: make([]*statement, 0),
        }
        assemblerlex.(*lexer).program = $$
    }
    | program stmt LF {
        fmt.Printf("* appendind stmt %v, stmt count %d\n", $2, len($$.statements))
        $$ = &program {
            statements: append($1.statements, $2),
        }
        assemblerlex.(*lexer).program = $$
    }

stmt: /* empty */ {
        fmt.Println("* comment or empty stmt")
        $$ = &statement{
            opcode: "comment",
        }
    }
    | lui_stmt { $$ = $1 }
    | auipc_stmt { $$ = $1 }
    | addi_stmt { $$ = $1 }
    | li_stmt { $$ = $1 }
    | add_stmt { $$ = $1 }
    | label_stmt { $$ = $1}
    | expr {
        fmt.Printf("* stmt expr %v\n", $$)
        $$ = &statement{
            opcode: "expr",
        }
    }

lui_stmt: LUI REGISTER COMMA NUMBER {
    fmt.Printf("* lui_stmt: %+v\n", $1)
    val, err := strconv.Atoi($4.lit)
    chkerr(err)
    $$ = &statement{
        opcode: $1.lit,
        op1: regs[$2.lit],
        op2: val,
    }
}

auipc_stmt: AUIPC REGISTER COMMA NUMBER {
    fmt.Printf("* auipc_stmt: %+v\n", $1)
    val, err := strconv.Atoi($4.lit)
    chkerr(err)
    $$ = &statement{
        opcode: $1.lit,
        op1: regs[$2.lit],
        op2: val,
    }
}

addi_stmt: ADDI REGISTER COMMA REGISTER COMMA NUMBER {
    fmt.Printf("* addi_stmt: %+v\n", $1)
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
    fmt.Printf("* li_stmt: %+v\n", $1)
    val, err := strconv.Atoi($4.lit)
    chkerr(err)
    $$ = &statement{
        opcode: $1.lit,
        op1: regs[$2.lit],
        op2: val,
    }
}

add_stmt: ADD REGISTER COMMA REGISTER COMMA REGISTER {
    fmt.Printf("* add_stmt: %+v\n", $1)
    $$ = &statement{
        opcode: $1.lit,
        op1: regs[$2.lit],
        op2: regs[$4.lit],
        op3: regs[$6.lit],
    }
}

label_stmt: IDENT COLON {
    fmt.Printf("* label_stmt: %+v\n", $1)
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
func checkerr(err error) {
    if err != nil {
        panic(err)
    }
}