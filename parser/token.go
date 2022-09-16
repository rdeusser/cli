package parser

import (
	"fmt"

	"github.com/rdeusser/cli/token"
)

type Token struct {
	Literal  string
	StartPos token.Pos
	EndPos   token.Pos
}

// Pos returns the start and end positions of a token.
func (t Token) Pos() string {
	return fmt.Sprintf("%s:%s", t.StartPos.String(), t.EndPos.String())
}

func (t Token) IsInvalid() bool {
	return t.Literal == "" && t.StartPos.String() == ""
}
