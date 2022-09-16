package token

import (
	"fmt"
)

// Pos is the position of a token.
type Pos struct {
	Column int
}

// NewPosition returns a position set at column 0.
func NewPosition() Pos {
	return Pos{
		Column: 1,
	}
}

// String returns the string form of a position.
func (p Pos) String() string {
	return fmt.Sprintf("%d", p.Column)
}
