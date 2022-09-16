// package multierror is inspired by github.com/hashicorp/go-multierror.
package multierror

import (
	"fmt"
	"strings"
)

type Error struct {
	errors []error
}

func (e *Error) Append(err error) {
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e *Error) ErrorOrNil() error {
	if e == nil {
		return nil
	}

	if len(e.errors) == 0 {
		return nil
	}

	return fmt.Errorf(e.Error())
}

func (e *Error) Error() string {
	if len(e.errors) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n", e.errors[0])
	}

	errors := make([]string, 0)
	for _, err := range e.errors {
		errors = append(errors, fmt.Sprintf("* %s", err))
	}

	return fmt.Sprintf("%d errors occurred:\n\t%s\n", len(errors), strings.Join(errors, "\n\t"))
}
