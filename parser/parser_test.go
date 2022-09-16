package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli/ast"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		testName string
		args     []string
		stmt     *ast.Statement
		want     string
	}{
		{
			"go run main.go",
			[]string{"/tmp/foo_bar/bar-baz/exe/main"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "/tmp/foo_bar/bar-baz/exe/main",
						Position: 0,
					},
				},
			},
			"/tmp/foo_bar/bar-baz/exe/main",
		},
		{
			"commands only",
			[]string{"kubectl", "get", "pods"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "kubectl",
						Position: 0,
					},
					{
						Name:     "get",
						Position: 1,
					},
					{
						Name:     "pods",
						Position: 2,
					},
				},
			},
			"kubectl get pods",
		},
		{
			"commands with flags",
			[]string{"kubectl", "get", "pods", "-o", "yaml"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "kubectl",
						Position: 0,
					},
					{
						Name:     "get",
						Position: 1,
					},
					{
						Name:     "pods",
						Position: 2,
					},
					{
						Name:     "-o",
						Position: 3,
					},
					{
						Name:     "yaml",
						Position: 4,
					},
				},
			},
			"kubectl get pods -o yaml",
		},
		{
			"commands with flags and arguments",
			[]string{"kubectl", "get", "pod", "-o", "yaml", "foo"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "kubectl",
						Position: 0,
					},
					{
						Name:     "get",
						Position: 1,
					},
					{
						Name:     "pod",
						Position: 2,
					},
					{
						Name:     "-o",
						Position: 3,
					},
					{
						Name:     "yaml",
						Position: 4,
					},
					{
						Name:     "foo",
						Position: 5,
					},
				},
			},
			"kubectl get pod -o yaml foo",
		},
		{
			"commands with flags and arguments in different order",
			[]string{"kubectl", "get", "pod", "foo", "-o", "yaml"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "kubectl",
						Position: 0,
					},
					{
						Name:     "get",
						Position: 1,
					},
					{
						Name:     "pod",
						Position: 2,
					},
					{
						Name:     "foo",
						Position: 3,
					},
					{
						Name:     "-o",
						Position: 4,
					},
					{
						Name:     "yaml",
						Position: 5,
					},
				},
			},
			"kubectl get pod foo -o yaml",
		},
		{
			"commands with flags and arguments and an unknown",
			[]string{"kubectl", "get", "pod", "foo", "-o", "yaml", "bar"},
			&ast.Statement{
				Arguments: []*ast.Argument{
					{
						Name:     "kubectl",
						Position: 0,
					},
					{
						Name:     "get",
						Position: 1,
					},
					{
						Name:     "pod",
						Position: 2,
					},
					{
						Name:     "foo",
						Position: 3,
					},
					{
						Name:     "-o",
						Position: 4,
					},
					{
						Name:     "yaml",
						Position: 5,
					},
					{
						Name:     "bar",
						Position: 6,
					},
				},
			},
			"kubectl get pod foo -o yaml bar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			p := New(tc.args)
			stmt := p.Parse()
			assert.NotNil(t, stmt)
			assert.Equal(t, tc.want, stmt.String(), "The arguments or the order of them didn't match what *ast.Statement output")
		})
	}
}
