package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAST(t *testing.T) {
	testCases := []struct {
		testName string
		node     Node
		want     string
	}{
		{
			"commands",
			&Statement{
				Arguments: []*Argument{
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
			`kubectl get pods`,
		},
		{
			"commands with flags",
			&Statement{
				Arguments: []*Argument{
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
			`kubectl get pods -o yaml`,
		},
		{
			"commands with arguments",
			&Statement{
				Arguments: []*Argument{
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
				},
			},
			`kubectl get pod foo`,
		},
		{
			"commands with flags and arguments",
			&Statement{
				Arguments: []*Argument{
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
			`kubectl get pod -o yaml foo`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.node.String())
		})
	}
}
