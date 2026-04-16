%{
package parser

import (
	"strconv"

	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

var parseErrors []string

%}

%union {
	program          *ast.Program
	statements       []ast.Statement
	statement        ast.Statement
	blockStatement   *ast.BlockStatement
	expression       ast.Expression
	expressions      []ast.Expression
	identifiers      []*ast.Identifier
	objPairs         map[string]ast.Expression
	token            token.Token
	literal          []rune
	intVal           int64
}

/* Tokens */
%token <token>  IDENTIFIER INT STRING
%token <token>  PLUS MINUS MULTIPLY DIVIDE MODULUS
%token <token>  EQUAL NOT_EQUAL GREATER_THAN LESS_THAN GREATER_THAN_OR_EQUAL LESS_THAN_OR_EQUAL
%token <token>  ASSIGNMENT COMMA SEMICOLON COLON DOT
%token <token>  LPAREN RPAREN LBRACKET RBRACKET LBRACE RBRACE
%token <token>  VAR FUNC RETURN IF ELSE NIL TRUE FALSE AND OR NOT

%type <program>         program
%type <statements>      statements
%type <statement>       statement
%type <blockStatement>  block
%type <expression>      expression
%type <expression>      primary
%type <expressions>     expressionList
%type <expressions>     arguments
%type <identifiers>     parameters
%type <objPairs>        objectPairs
%type <objPairs>        objectPairsList

/* Operator precedence and associativity */
%left OR
%left AND
%left EQUAL NOT_EQUAL
%left GREATER_THAN LESS_THAN GREATER_THAN_OR_EQUAL LESS_THAN_OR_EQUAL
%left PLUS MINUS
%left MULTIPLY DIVIDE MODULUS
%right UNARY_MINUS UNARY_NOT
%left DOT
%left LPAREN RPAREN LBRACKET RBRACKET

%%

program
	: statements
	{
		$$ = &ast.Program{Statements: $1}
		yylex.(*YaccLexer).program = $$
	}
	| /* empty */
	{
		$$ = &ast.Program{Statements: []ast.Statement{}}
		yylex.(*YaccLexer).program = $$
	}
	;

statements
	: statement
	{
		if $1 != nil {
			$$ = []ast.Statement{$1}
		} else {
			$$ = []ast.Statement{}
		}
	}
	| statements statement
	{
		if $2 != nil {
			$$ = append($1, $2)
		} else {
			$$ = $1
		}
	}
	;

statement
	: VAR IDENTIFIER ASSIGNMENT expression optSemicolon
	{
		$$ = &ast.VarStatement{
			Token: $1,
			Name: &ast.Identifier{
				Token: $2,
				Value: $2.Literal,
			},
			Value: $4,
		}
	}
	| RETURN expression optSemicolon
	{
		$$ = &ast.ReturnStatement{
			Token:       $1,
			ReturnValue: $2,
		}
	}
	| expression optSemicolon
	{
		stmt := &ast.ExpressionStatement{Expression: $1}
		if expr, ok := $1.(*ast.Identifier); ok {
			stmt.Token = expr.Token
		} else if lit, ok := $1.(*ast.IntegerLiteral); ok {
			stmt.Token = lit.Token
		} else if str, ok := $1.(*ast.String); ok {
			stmt.Token = str.Token
		} else if prefixExpr, ok := $1.(*ast.PrefixExpression); ok {
			stmt.Token = prefixExpr.Token
		} else if infixExpr, ok := $1.(*ast.InfixExpression); ok {
			stmt.Token = infixExpr.Token
		} else if callExpr, ok := $1.(*ast.CallExpression); ok {
			stmt.Token = callExpr.Token
		} else if indexExpr, ok := $1.(*ast.IndexExpression); ok {
			stmt.Token = indexExpr.Token
		} else if listExpr, ok := $1.(*ast.List); ok {
			stmt.Token = listExpr.Token
		} else if boolExpr, ok := $1.(*ast.Boolean); ok {
			stmt.Token = boolExpr.Token
		} else if nilExpr, ok := $1.(*ast.Nil); ok {
			stmt.Token = nilExpr.Token
		}
		$$ = stmt
	}
	| error
	{
		// Let yacc's default error handling record the error
		$$ = nil
	}
	;

optSemicolon
	: SEMICOLON
	| /* empty */
	;

block
	: LBRACE statements RBRACE
	{
		$$ = &ast.BlockStatement{
			Token:      $1,
			Statements: $2,
		}
	}
	| LBRACE RBRACE
	{
		$$ = &ast.BlockStatement{
			Token:      $1,
			Statements: []ast.Statement{},
		}
	}
	;

expression
	: primary
	| expression PLUS expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression MINUS expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression MULTIPLY expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression DIVIDE expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression MODULUS expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression EQUAL expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression NOT_EQUAL expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression GREATER_THAN expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression LESS_THAN expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression GREATER_THAN_OR_EQUAL expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression LESS_THAN_OR_EQUAL expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression AND expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| expression OR expression
	{
		$$ = &ast.InfixExpression{
			Token:    $2,
			Left:     $1,
			Operator: string($2.Literal),
			Right:    $3,
		}
	}
	| MINUS expression %prec UNARY_MINUS
	{
		$$ = &ast.PrefixExpression{
			Token:    $1,
			Operator: string($1.Literal),
			Right:    $2,
		}
	}
	| NOT expression %prec UNARY_NOT
	{
		$$ = &ast.PrefixExpression{
			Token:    $1,
			Operator: string($1.Literal),
			Right:    $2,
		}
	}
	| expression LBRACKET expression RBRACKET
	{
		$$ = &ast.IndexExpression{
			Token: $2,
			List:  $1,
			Index: $3,
		}
	}
	| expression DOT IDENTIFIER
	{
		$$ = &ast.PropertyAccess{
			Token:    $2,
			Object:   $1,
			Property: string($3.Literal),
		}
	}
	| expression LPAREN arguments RPAREN
	{
		$$ = &ast.CallExpression{
			Token:     $2,
			Function:  $1,
			Arguments: $3,
		}
	}
	;

primary
	: IDENTIFIER
	{
		$$ = &ast.Identifier{
			Token: $1,
			Value: $1.Literal,
		}
	}
	| INT
	{
		val, _ := strconv.ParseInt(string($1.Literal), 0, 64)
		$$ = &ast.IntegerLiteral{
			Token: $1,
			Value: val,
		}
	}
	| STRING
	{
		$$ = &ast.String{
			Token: $1,
			Value: $1.Literal,
		}
	}
	| TRUE
	{
		$$ = &ast.Boolean{
			Token: $1,
			Value: true,
		}
	}
	| FALSE
	{
		$$ = &ast.Boolean{
			Token: $1,
			Value: false,
		}
	}
	| NIL
	{
		$$ = &ast.Nil{Token: $1}
	}
	| LBRACKET expressionList RBRACKET
	{
		$$ = &ast.List{
			Token: $1,
			Items: $2,
		}
	}
	| LBRACKET RBRACKET
	{
		$$ = &ast.List{
			Token: $1,
			Items: []ast.Expression{},
		}
	}
	| LBRACE objectPairs RBRACE
	{
		$$ = &ast.ObjectLiteral{
			Token: $1,
			Pairs: $2,
		}
	}
	| LBRACE RBRACE
	{
		$$ = &ast.ObjectLiteral{
			Token: $1,
			Pairs: make(map[string]ast.Expression),
		}
	}
	| LPAREN expression RPAREN
	{
		$$ = $2
	}
	| IF LPAREN expression RPAREN block
	{
		$$ = &ast.IfExpression{
			Token:       $1,
			Condition:   $3,
			Consequence: $5,
		}
	}
	| IF LPAREN expression RPAREN block ELSE block
	{
		$$ = &ast.IfExpression{
			Token:       $1,
			Condition:   $3,
			Consequence: $5,
			Alternative: $7,
		}
	}
	| FUNC LPAREN parameters RPAREN block
	{
		$$ = &ast.FuncExpression{
			Token:  $1,
			Params: $3,
			Body:   $5,
		}
	}
	;

expressionList
	: expression
	{
		$$ = []ast.Expression{$1}
	}
	| expressionList COMMA expression
	{
		$$ = append($1, $3)
	}
	;

arguments
	: expression
	{
		$$ = []ast.Expression{$1}
	}
	| arguments COMMA expression
	{
		$$ = append($1, $3)
	}
	| /* empty */
	{
		$$ = []ast.Expression{}
	}
	;

parameters
	: IDENTIFIER
	{
		$$ = []*ast.Identifier{
			{
				Token: $1,
				Value: $1.Literal,
			},
		}
	}
	| parameters COMMA IDENTIFIER
	{
		$$ = append($1, &ast.Identifier{
			Token: $3,
			Value: $3.Literal,
		})
	}
	| /* empty */
	{
		$$ = []*ast.Identifier{}
	}
	;

objectPairs
	: objectPairsList
	{
		$$ = $1
	}
	;

objectPairsList
	: IDENTIFIER COLON expression
	{
		$$ = make(map[string]ast.Expression)
		$$[string($1.Literal)] = $3
	}
	| objectPairsList COMMA IDENTIFIER COLON expression
	{
		$1[string($3.Literal)] = $5
		$$ = $1
	}
	;

%%

type YaccLexer struct {
	impl    *lexer.LexerImpl
	program *ast.Program
}

func (l *YaccLexer) Error(s string) {
	parseErrors = append(parseErrors, s)
}

func (l *YaccLexer) Lex(lval *yySymType) int {
	tkn := l.impl.NextToken()
	
	if tkn.Type == token.EOF {
		return 0
	}
	
	lval.token = tkn
	
	switch tkn.Type {
	case token.IDENTIFIER:
		return IDENTIFIER
	case token.INT:
		return INT
	case token.STRING:
		return STRING
	case token.ASSIGNMENT:
		return ASSIGNMENT
	case token.PLUS:
		return PLUS
	case token.MINUS:
		return MINUS
	case token.MULTIPLY:
		return MULTIPLY
	case token.DIVIDE:
		return DIVIDE
	case token.MODULUS:
		return MODULUS
	case token.EQUAL:
		return EQUAL
	case token.NOT_EQUAL:
		return NOT_EQUAL
	case token.GREATER_THAN:
		return GREATER_THAN
	case token.LESS_THAN:
		return LESS_THAN
	case token.GREATER_THAN_OR_EQUAL:
		return GREATER_THAN_OR_EQUAL
	case token.LESS_THAN_OR_EQUAL:
		return LESS_THAN_OR_EQUAL
	case token.COMMA:
		return COMMA
	case token.SEMICOLON:
		return SEMICOLON
	case token.COLON:
		return COLON
	case token.DOT:
		return DOT
	case token.LPAREN:
		return LPAREN
	case token.RPAREN:
		return RPAREN
	case token.LBRACKET:
		return LBRACKET
	case token.RBRACKET:
		return RBRACKET
	case token.LBRACE:
		return LBRACE
	case token.RBRACE:
		return RBRACE
	case token.VAR:
		return VAR
	case token.FUNC:
		return FUNC
	case token.RETURN:
		return RETURN
	case token.IF:
		return IF
	case token.ELSE:
		return ELSE
	case token.NIL:
		return NIL
	case token.TRUE:
		return TRUE
	case token.FALSE:
		return FALSE
	case token.AND:
		return AND
	case token.OR:
		return OR
	case token.NOT:
		return NOT
	}
	
	return int(tkn.Type)
}
