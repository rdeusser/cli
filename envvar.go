package cli

import (
	"os"
)

type EnvVar[T Value] struct {
	// Name of the environment variable.
	Name string

	// Layout is the layout to use if the environment variable should be parsed
	// as a time.Time value.
	Layout string
}

func (e *EnvVar[T]) Lookup() (T, error) {
	var result T

	if e.Name == "" {
		return result, ErrEnvVarMustHaveName
	}

	if env, ok := os.LookupEnv(e.Name); ok {
		return parseValue[T](env, 0, e.Layout)
	}

	return result, nil
}
