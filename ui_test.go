package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	SetOutput(&buf)
	f()
	return buf.String()
}

func TestOutput(t *testing.T) {
	SetProjectName("foo")

	t.Run("Output", func(t *testing.T) {
		output := captureOutput(func() {
			Output("hello world")
		})

		expected := "hello world\n"

		assert.Equal(t, expected, output)
	})

	t.Run("Info", func(t *testing.T) {
		output := captureOutput(func() {
			Info("hello world")
		})

		expected := fmt.Sprintf("%s%s\n", colorize(ColorGreen, "[INFO][foo]: "), "hello world")

		assert.Equal(t, expected, output)
	})

	t.Run("Warn", func(t *testing.T) {
		output := captureOutput(func() {
			Warn("hello world")
		})

		expected := fmt.Sprintf("%s%s\n", colorize(ColorGreen, "[WARN][foo]: "), "hello world")

		assert.Equal(t, expected, output)
	})

	t.Run("Error", func(t *testing.T) {
		output := captureOutput(func() {
			Error("hello world")
		})

		expected := fmt.Sprintf("%s%s\n", colorize(ColorGreen, "[ERROR][foo]: "), "hello world")

		assert.Equal(t, expected, output)
	})

	t.Run("SetProjectName", func(t *testing.T) {
		SetProjectName("bar")

		output := captureOutput(func() {
			Info("hello world")
		})

		expected := fmt.Sprintf("%s%s\n", colorize(ColorGreen, "[INFO][bar]: "), "hello world")

		assert.Equal(t, expected, output)
	})
}
