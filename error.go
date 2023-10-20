package serrors

import (
	"errors"
	"fmt"
)

type serror struct {
	err     error
	keyvals []any
}

func (e *serror) Error() string {
	if e.err == nil {
		return "nil error"
	}

	return e.err.Error()
}

func (e *serror) KeyVals() []any {
	return e.keyvals
}

func (e *serror) Unwrap() error {
	return errors.Unwrap(e.err)
}

func New(msg string, args ...any) error {
	return &serror{
		err:     fmt.Errorf(msg),
		keyvals: args,
	}
}

func Wrap(err error, msg string, args ...any) error {
	return &serror{
		err:     fmt.Errorf("%s: %w", msg, err),
		keyvals: args,
	}
}

type keyvaluer interface {
	KeyVals() []any
}

func KeyVals(err error) []any {
	var keyvals []any

	for err != nil {
		if u, ok := err.(keyvaluer); ok {
			keyvals = append(keyvals, u.KeyVals()...)
		}

		err = errors.Unwrap(err)
	}

	return keyvals
}
