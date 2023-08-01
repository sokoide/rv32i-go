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
%type<stmt> stmt lui_stmt auipc_stmt jal_stmt jalr_stmt ret_stmt
%type<stmt> lw_stmt lbu_stmt sb_stmt sh_stmt sw_stmt
%type<stmt> addi_stmt li_stmt sltiu_stmt andi_stmt srli_stmt seqz_stmt
%type<stmt> add_stmt sub_stmt
%type<stmt> label_stmt
%type<expr> expr
%token<tok> LF COLON COMMA LP RP NUMBER IDENT
%token<tok> REGISTER AUIPC LUI JAL JALR RET LW LBU SB SH SW ADDI LI SLTIU SEQZ ANDI SRLI
%token<tok> ADD SUB


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
    | jal_stmt { $$ = $1 }
    | jalr_stmt { $$ = $1 }
    | ret_stmt { $$ = $1 }
    | lw_stmt { $$ = $1 }
    | lbu_stmt { $$ = $1 }
    | sb_stmt { $$ = $1 }
    | sh_stmt { $$ = $1 }
    | sw_stmt { $$ = $1 }
    | addi_stmt { $$ = $1 }
    | li_stmt { $$ = $1 }
    | sltiu_stmt { $$ = $1 }
    | seqz_stmt { $$ = $1 }
    | andi_stmt { $$ = $1 }
    | srli_stmt { $$ = $1 }
    | add_stmt { $$ = $1 }
    | sub_stmt { $$ = $1 }
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
            op1: 1, // if rd is omitted, defaults to x1
            str1: $2.lit,
        }
    }

jalr_stmt: JALR REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* jalr_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
			op3: regs[$6.lit],
        }
    }
    | JALR NUMBER LP REGISTER RP {
        log.Debugf("* jalr_stmt: %+v", $1)
        val, err := strconv.Atoi($2.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: 1,
            op2: val,
			op3: regs[$4.lit],
        }
    }

ret_stmt: RET {
        log.Debugf("* ret_stmt")
        $$ = &statement{
            opcode: "jalr",
            op1: 0,
            op2: 0,
            op3: 1,
        }
    }

lw_stmt: LW REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lw_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
            op3: regs[$6.lit],
        }
    }

lbu_stmt: LBU REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lbu_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
            op3: regs[$6.lit],
        }
    }

sb_stmt: SB REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sb_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
            op3: regs[$6.lit],
        }
    }

sh_stmt: SH REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sh_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
            op3: regs[$6.lit],
        }
    }

sw_stmt: SW REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sw_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: val,
            op3: regs[$6.lit],
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

sltiu_stmt: SLTIU REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* sltiu_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: val,
        }
    }

seqz_stmt: SEQZ REGISTER COMMA REGISTER {
        log.Debugf("* seqz_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
        }
    }

andi_stmt: ANDI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* andi_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: val,
        }
}


srli_stmt: SRLI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* srli_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: val,
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

sub_stmt: SUB REGISTER COMMA REGISTER COMMA REGISTER {
        log.Debugf("* sub_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: regs[$2.lit],
            op2: regs[$4.lit],
            op3: regs[$6.lit],
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
