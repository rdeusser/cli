package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	args := []string{"kubectl", "get", "pods"}
	want := []string{"pods"}
	buf := Reduce(args, func(a string) bool {
		return a == "pods"
	})

	assert.Equal(t, want, buf)
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		i    int
		j    int
		want []string
	}{
		{"remove namespace flag", []string{"kubectl", "get", "pods", "-n", "kube-system"}, 3, 5, []string{"kubectl", "get", "pods"}},
		{"remove argument", []string{"kubectl", "get", "pods", "foo", "-n", "kube-system"}, 3, 4, []string{"kubectl", "get", "pods", "-n", "kube-system"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := Remove(tc.args, tc.i, tc.j)

			assert.Equal(t, tc.want, args)
		})
	}
}

func TestRemoveForLoop(t *testing.T) {
	args := []string{"kubectl", "get", "pods", "-n", "kube-system"}
	want := []string{"kubectl", "get", "pods"}
	for i, arg := range args {
		if arg == "-n" {
			args = Remove(args, i, i+2)
		}
	}

	assert.Equal(t, want, args)
}
