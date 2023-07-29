%{
package main
%}

%union{
    expr    expression
    tok     token
}

%type<expr> program expr
%token<tok> LF NUMBER IDENT

%left '+' '-'
%left '*' '/'

%start program

%%
/* program: line {
        assemblerlex.(*lexer).program = $$
    }
    | program line {
        assemblerlex.(*lexer).program =$$
    }

line: label LF {
        assemblerlex.(*lexer).line =$$
    }

label: IDENT ':' {
        assemblerlex.(*lexer).label =$$
    } */

program: expr {
		$$ = $1
        assemblerlex.(*lexer).program = $$
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
