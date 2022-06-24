package errors

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type withMessage struct {
	err error
	msg string
}

// Error returns an error string with the message and cause concatenated.
func (w *withMessage) Error() string { return w.msg + ": " + w.err.Error() }

// Is implements the inline interface (`err.(interface{ Is(error) bool })`) that
// the standard libary errors.Is looks for in errors.
func (w *withMessage) Is(target error) bool { return errors.Is(w.err, target) }

// As implements the inline interface (`err.(interface{ As(any) bool })`) that
// the standard libary errors.As looks for in errors.
func (w *withMessage) As(target interface{}) bool { return errors.As(w.err, target) }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMessage) Unwrap() error { return w.err }

// Format provides a method for fmt.Sprint(f) or Fprint(f) to generate output
// with.
func (w *withMessage) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%+v\n", w.err)
			io.WriteString(f, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(f, w.Error())
	}
}

// As wraps the standard library's As function to avoid name collisions.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is is functionally the same as errors.Is, but instead looks through the
// entire chain of errors and strings to find target (if it can).
//
// Target does not need to implement the error type. Some packages will just
// return an `fmt.Errorf` or `errors.New` and that makes it impossible to check
// normally. Here you can just pas
func Is(err error, target interface{}) bool {
	if target == nil {
		return err == target
	}

	// First we work our way through the `withMessage` error type, checking
	// to see if our error matches target along the way.
	if terr, ok := target.(*withMessage); ok {
		if errors.Is(err, terr) {
			return true
		}

		return Is(err, terr.err)
	}

	// If we've gone through all the `withMessage` error types, check to see
	// if we have regular error types now that may be returned as
	// `fmt.Errorf` or `errors.New`.
	if terr, ok := target.(error); ok {
		if errors.Is(err, terr) {
			return true
		}

		// If the error still doesn't match, attempt to unwrap the error
		// and check that.
		if uerr := errors.Unwrap(terr); uerr != nil {
			return Is(err, uerr)
		}

		// If target doesn't match our error and we can't unwrap it
		// anymore return the error string and check that.
		return Is(err, terr.Error())
	}

	// Finally, if the error string matches the target string return
	// true. Otherwise there definitely can't be a match.
	if s, ok := target.(string); ok {
		if strings.TrimSpace(err.Error()) == strings.TrimSpace(s) {
			return true
		}
	}

	return false
}

// Unwrap wraps the standard library's Unwrap function to avoid name collisions
// in downstream packages.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Wrap wraps an error and a corresponding message. It helps with discovering
// where an error first occurred and the chain of events it caused.
func Wrap(err error, msg string) error {
	return &withMessage{
		err: err,
		msg: msg,
	}
}

// Wrapf wraps an error and a corresponding message. It helps with discovering
// where an error first occurred and the chain of events it caused.
func Wrapf(err error, format string, args ...interface{}) error {
	return &withMessage{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}
