package terrors

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"

	"github.com/go-faster/errors"
)

type noWrapper struct {
	error
}

func (e noWrapper) FormatError(p errors.Printer) (next error) {
	if f, ok := e.error.(errors.Formatter); ok {
		return f.FormatError(p)
	}
	p.Print(e.error)
	return nil
}

type wrapError struct {
	msg   string
	err   error
	frame errors.Frame
}

var _ Framer = (*wrapError)(nil)

func (e *wrapError) Root() error {
	return e.err
}

func (e *wrapError) Frame() errors.Frame {
	return e.frame
}

func (e *wrapError) Info() []any {
	return []any{e.msg}
}

func (e *wrapError) Error() string {
	return fmt.Sprint(e)
}

func (e *wrapError) Format(s fmt.State, v rune) { errors.FormatError(e, s, v) }

func (e *wrapError) FormatError(p errors.Printer) (next error) {
	p.Print(e.msg)
	e.frame.Format(p)
	return e.err
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// Wrap error with message and caller.
func Wrap(err error, message string) error {
	frame := errors.Frame{}
	if errors.Trace() {
		frame = errors.Caller(1)
	}
	return &wrapError{msg: message, err: err, frame: frame}
}

// Wrapf wraps error with formatted message and caller.
func Wrapf(err error, format string, a ...interface{}) error {
	frame := errors.Frame{}
	if errors.Trace() {
		frame = errors.Caller(1)
	}
	msg := fmt.Sprintf(format, a...)
	return &wrapError{msg: msg, err: err, frame: frame}
}
