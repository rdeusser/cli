package parser

import (
	"strings"

	"github.com/rdeusser/cli/ast"
	"github.com/rdeusser/cli/token"
)

type Parser struct {
	args     []string
	joined   string
	pos      int
	curToken Token
}

func New(args []string) *Parser {
	return &Parser{
		args:     args,
		joined:   strings.Join(args, " "),
		pos:      -1,
		curToken: Token{},
	}
}

func (p *Parser) Parse() *ast.Statement {
	p.nextToken()

	stmt := &ast.Statement{
		Arguments: make([]*ast.Argument, 0),
		From:      p.curToken.StartPos,
	}

	for i := range p.args {
		stmt.Arguments = append(stmt.Arguments, p.parseArgument(i))
	}

	stmt.To = p.curToken.EndPos

	return stmt
}

func (p *Parser) EndOfArgs() bool {
	return p.pos >= len(p.args)
}

func (p *Parser) parseArgument(pos int) *ast.Argument {
	arg := &ast.Argument{
		From: p.curToken.StartPos,
	}

	arg.Name = p.curToken.Literal
	arg.Position = pos
	arg.To = p.curToken.EndPos

	p.nextToken()

	return arg
}

func (p *Parser) nextToken() {
	p.pos++

	if p.EndOfArgs() {
		return
	}

	arg := p.args[p.pos]
	start := strings.Index(p.joined, arg)
	end := start + len(arg) + 1
	tok := Token{
		Literal: arg,
		StartPos: token.Pos{
			Column: start,
		},
		EndPos: token.Pos{
			Column: end,
		},
	}

	if tok.StartPos.Column == 0 {
		tok.StartPos.Column = 1
	}

	p.curToken = tok
}
