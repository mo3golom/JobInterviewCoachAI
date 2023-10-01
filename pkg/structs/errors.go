package structs

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
)

type (
	BaseError struct {
		error
	}

	RecoveredErr struct {
		Stack []byte
		Cause error
	}
)

func NewBaseError(err error) BaseError {
	return BaseError{err}
}

func (b BaseError) Unwrap() error {
	err := b.error
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return nil
}

func (e *RecoveredErr) Error() string {
	return e.Cause.Error()
}

func WithRecover(block func() error) (err error) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			err = &RecoveredErr{
				Stack: debug.Stack(),
				Cause: fmt.Errorf("recovered from: %s", recovered),
			}
		}
	}()

	return block()
}
