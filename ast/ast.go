package ast

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rdeusser/cli/internal/join"
	"github.com/rdeusser/cli/token"
)

type Node interface {
	fmt.Stringer
	Pos() (start, end int)
}

type Statement struct {
	Arguments []*Argument
	From      token.Pos
	To        token.Pos
}

func (n *Statement) Lookup(s string) *Argument {
	for _, arg := range n.Arguments {
		if arg.Name == s {
			return arg
		}
	}

	return nil
}

func (n *Statement) String() string {
	args := make([]string, 0)

	sort.Sort(SortArgumentsByPosition(n.Arguments))

	for _, arg := range n.Arguments {
		args = append(args, arg.String())
	}

	return join.Args(args)
}

func (n *Statement) Pos() (start, end int) {
	start = n.From.Column
	end = n.To.Column
	return start, end
}

type Argument struct {
	Name     string
	Position int
	From     token.Pos
	To       token.Pos
}

func (n *Argument) String() string {
	return strings.TrimSpace(n.Name)
}

func (n *Argument) Pos() (start, end int) {
	start = n.From.Column
	end = n.To.Column
	return start, end
}

// SortArgumentsByPosition sorts args by position.
type SortArgumentsByPosition []*Argument

func (n SortArgumentsByPosition) Len() int           { return len(n) }
func (n SortArgumentsByPosition) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortArgumentsByPosition) Less(i, j int) bool { return n[i].Position < n[j].Position }
