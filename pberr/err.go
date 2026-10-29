// Package pberr provides helpers for working with errors.
package pberr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/peterbourgon/pbgo/pbslices"
)

// Annotatef prefix-annotates err with the provided Printf-style format and
// args. Unlike fmt.Errorf, if err is nil, Annotatef returns nil.
func Annotatef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	if format == "" {
		return err
	}

	return fmt.Errorf(format+": %w", append(args, err)...)
}

// Squelch returns nil if errors.Is(err, squelched) returns true for any of the
// provided squelched values, and err otherwise.
func Squelch(err error, squelched ...error) error {
	for _, s := range squelched {
		if errors.Is(err, s) {
			return nil
		}
	}
	return err
}

// Join is like [errors.Join] but concatenates the error messages with ';'
// instead of newlines.
func Join(errs ...error) error {
	if es := pbslices.Filter(errs, func(err error) bool {
		return err != nil
	}); len(es) > 0 {
		return &joinError{
			errs: es,
		}

	}
	return nil
}

// joinError is the error type returned by Join.
type joinError struct {
	errs []error
}

// Error implements the error interface.
func (e *joinError) Error() string {
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	default:
		return strings.Join(pbslices.Transform(e.errs, error.Error), "; ")
	}
}

// Unwrap returns the errors passed to Join.
func (e *joinError) Unwrap() []error {
	return e.errs
}
