%{
package main

import "fmt"

%}

%union{
    program program
    stmt    statement
    expr    expression
    tok     token
}

%type<program> program
%type<stmt> stmt label_stmt li_stmt lui_stmt
%type<expr> expr
%token<tok> LF COLON COMMA NUMBER IDENT
%token<tok> REGISTER LI LUI

%left '+' '-'
%left '*' '/'

%start program

%%
program: /* empty */ {
    fmt.Println("* empty program")
        $$ = program{
            statements: make([]statement, 0),
        }
        assemblerlex.(*lexer).program = &$$
    }
    | program stmt LF {
        // $$.statements = append($1.statements, $2)
        fmt.Printf("* appendind stmt %v, stmt count %d\n", $2, len($$.statements))
        $$ = program {
            statements: append($1.statements, $2),
        }
        assemblerlex.(*lexer).program = &$$
    }

stmt: /* empty */ {
        fmt.Println("* comment or empty stmt")
        $$ = nil
    }
    | label_stmt { $$ = $1}
    | li_stmt { $$ = $1 }
    | lui_stmt { $$ = $1 }
    | expr {
        fmt.Printf("* stmt expr %v\n", $$)
        $$ = $1
    }

label_stmt: IDENT COLON {
    fmt.Printf("* label_stmt: %+v\n", $1)
        $$ = $1
}

li_stmt: LI REGISTER COMMA NUMBER {
    fmt.Printf("* li_stmt: %+v\n", $1)
        $$ = $1
}

lui_stmt: LUI REGISTER COMMA NUMBER {
    fmt.Printf("* lui_stmt: %+v\n", $1)
        $$ = $1
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