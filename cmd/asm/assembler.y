%{
package main

import "fmt"
%}

%union{
    stmt    statement
    expr    expression
    tok     token
}

%type<stmt> program stmt label_stmt li_stmt
%type<expr> expr
%token<tok> LF COLON COMMA NUMBER IDENT
%token<tok> REGISTER LI LUI

%left '+' '-'
%left '*' '/'

%start program

%%
program: /* empty */ {
        $$ = nil
    }
    | program stmt LF

stmt: /* empty */ {
    fmt.Println("* comment or empty stmt")
    }
    | label_stmt { }
    | li_stmt {}
    | lui_stmt {}
    | expr {
        $$ = $1
        // assemblerlex.(*lexer).program = $$
        // TODO:
        assemblerlex.(*lexer).program = nil
        fmt.Printf("* stmt expr %v\n", $$)
    }

label_stmt: IDENT COLON {
    fmt.Printf("* label_stmt: %+v\n", $1)
}

li_stmt: LI REGISTER COMMA NUMBER {
    fmt.Printf("* li_stmt: %+v\n", $1)
}

lui_stmt: LUI REGISTER COMMA NUMBER {
    fmt.Printf("* lui_stmt: %+v\n", $1)
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
