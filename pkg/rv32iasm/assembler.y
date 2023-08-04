%{
package rv32iasm

import (
    "strconv"
    "github.com/sokoide/rv32i-go/pkg/rv32i"

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
%type<stmt> beq_stmt bne_stmt blt_stmt bge_stmt bltu_stmt bgeu_stmt
%type<stmt> lb_stmt lh_stmt lw_stmt lbu_stmt lhu_stmt sb_stmt sh_stmt sw_stmt
%type<stmt> addi_stmt li_stmt slti_stmt sltiu_stmt seqz_stmt xori_stmt ori_stmt andi_stmt slli_stmt srli_stmt srai_stmt
%type<stmt> add_stmt sub_stmt
%type<stmt> label_stmt
%type<expr> expr
%token<tok> LF COLON COMMA LP RP NUMBER IDENT
%token<tok> REGISTER AUIPC LUI JAL JALR RET
%token<tok> BEQ BNE BLT BGE BLTU BGEU
%token<tok> LB LH LW LBU LHU SB SH SW
%token<tok> ADDI LI SLTI SLTIU SEQZ XORI ORI ANDI SLLI SRLI SRAI
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
    | beq_stmt { $$ = $1 }
    | bne_stmt { $$ = $1 }
    | blt_stmt { $$ = $1 }
    | bge_stmt { $$ = $1 }
    | bltu_stmt { $$ = $1 }
    | bgeu_stmt { $$ = $1 }
    | lb_stmt { $$ = $1 }
    | lh_stmt { $$ = $1 }
    | lw_stmt { $$ = $1 }
    | lbu_stmt { $$ = $1 }
    | lhu_stmt { $$ = $1 }
    | sb_stmt { $$ = $1 }
    | sh_stmt { $$ = $1 }
    | sw_stmt { $$ = $1 }
    | addi_stmt { $$ = $1 }
    | li_stmt { $$ = $1 }
    | slti_stmt { $$ = $1 }
    | sltiu_stmt { $$ = $1 }
    | seqz_stmt { $$ = $1 }
    | xori_stmt { $$ = $1 }
    | ori_stmt { $$ = $1 }
    | andi_stmt { $$ = $1 }
    | slli_stmt { $$ = $1 }
    | srli_stmt { $$ = $1 }
    | srai_stmt { $$ = $1 }
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
            op1: rv32i.Regs[$2.lit],
            op2: val,
        }
    }

auipc_stmt: AUIPC REGISTER COMMA NUMBER {
        log.Debugf("* auipc_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
        }
    }

jal_stmt: JAL REGISTER COMMA NUMBER {
        log.Debugf("* jal_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
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
            op1: rv32i.Regs[$2.lit],
            op2: val,
			op3: rv32i.Regs[$6.lit],
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
			op3: rv32i.Regs[$4.lit],
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

beq_stmt: BEQ REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* beq_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "beq",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

bne_stmt: BNE REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* bne_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "bne",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

blt_stmt: BLT REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* blt_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "blt",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

bge_stmt: BGE REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* bge_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "bge",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

bltu_stmt: BLTU REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* bltu_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "bltu",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

bgeu_stmt: BGEU REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* bgeu_stmt")
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: "bgeu",
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

lb_stmt: LB REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lb_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

lh_stmt: LH REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lh_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

lw_stmt: LW REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lw_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

lbu_stmt: LBU REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lbu_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

lhu_stmt: LHU REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* lhu_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

sb_stmt: SB REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sb_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

sh_stmt: SH REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sh_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

sw_stmt: SW REGISTER COMMA NUMBER LP REGISTER RP {
        log.Debugf("* sw_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
            op3: rv32i.Regs[$6.lit],
        }
    }

addi_stmt: ADDI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* addi_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

li_stmt: LI REGISTER COMMA NUMBER {
        log.Debugf("* li_stmt: %+v", $1)
        val, err := strconv.Atoi($4.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: val,
        }
    }

slti_stmt: SLTI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* slti_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

sltiu_stmt: SLTIU REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* sltiu_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

seqz_stmt: SEQZ REGISTER COMMA REGISTER {
        log.Debugf("* seqz_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
        }
    }

xori_stmt: XORI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* xori_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

ori_stmt: ORI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* ori_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

andi_stmt: ANDI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* andi_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
}

slli_stmt: SLLI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* slli_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

srli_stmt: SRLI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* srli_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

srai_stmt: SRAI REGISTER COMMA REGISTER COMMA NUMBER {
        log.Debugf("* srai_stmt: %+v", $1)
        val, err := strconv.Atoi($6.lit)
        chkerr(err)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: val,
        }
    }

add_stmt: ADD REGISTER COMMA REGISTER COMMA REGISTER {
        log.Debugf("* add_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: rv32i.Regs[$6.lit],
        }
    }

sub_stmt: SUB REGISTER COMMA REGISTER COMMA REGISTER {
        log.Debugf("* sub_stmt: %+v", $1)
        $$ = &statement{
            opcode: $1.lit,
            op1: rv32i.Regs[$2.lit],
            op2: rv32i.Regs[$4.lit],
            op3: rv32i.Regs[$6.lit],
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
