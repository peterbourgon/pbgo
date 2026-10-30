// Package pberr provides helpers for working with errors.
package pberr

import (
	"errors"
	"fmt"
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
