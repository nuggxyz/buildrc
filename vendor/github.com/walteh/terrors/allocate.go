package terrors

import (
	"fmt"

	"github.com/go-faster/errors"
)

type TracableError interface {
	error
	Trace(...any) error
	Child(string) TracableError
}

var _ TracableError = (*allocError)(nil)

var _ Framer = (*allocError)(nil)

// errorString is a trivial implementation of error.
type allocError struct {
	s     string
	root  error
	frame *errors.Frame
	info  []any
	s2    string
}

func New(str string) TracableError {
	return &allocError{s: str, root: nil, frame: nil, info: nil}
}

// New returns an error that formats as the given text.
//
// The returned error contains a Frame set to the caller's location and
// implements Formatter to show this information when printed with details.
func NewTraced(text string) error {
	if !errors.Trace() {
		return &allocError{text, nil, nil, nil, ""}
	}
	ofs := errors.Caller(1)
	return &allocError{text, nil, &ofs, nil, ""}
}

func (e *allocError) SetFrame(callerOffset int) {
	if !errors.Trace() {
		return
	}
	ofs := errors.Caller(callerOffset)
	e.frame = &ofs
}

func (e *allocError) Error() string { return e.s }

func (e *allocError) Format(s fmt.State, v rune) { errors.FormatError(e, s, v) }

func (e *allocError) FormatError(p errors.Printer) (next error) {
	if e.s2 != "" {
		p.Print(fmt.Sprintf("%s -> %s", e.s, e.s2))
	} else {
		p.Print(e.s)
	}
	if e.frame != nil {
		e.frame.Format(p)
	}
	return nil
}

func (e *allocError) Frame() errors.Frame {
	if e.frame == nil {
		return errors.Frame{}
	}
	return *e.frame
}

func (e *allocError) Root() error {
	return e.root
}

func (e *allocError) Info() []any {
	return e.info
}

func (g *allocError) Trace(info ...any) error {
	e := *g
	e.SetFrame(2)
	for _, v := range info {
		if f, ok := v.(error); ok {
			e.root = f
			break
		}
	}

	for _, v := range info {
		if f, ok := v.(string); ok {
			e.s2 = f
			break
		}
	}

	// if len(info) != 0 {
	// 	e.info = info
	// 	if len(info) >= 1 {
	// 		if f, ok := info[0].(error); ok {
	// 			e.root = f
	// 			e.info = info[1:]
	// 		} else if s, ok := info[0].(string); ok {
	// 			e.root = stderrors.New(s)
	// 			e.info = info[1:]
	// 		}
	// 	}
	// }
	return &e
}

func (e *allocError) Child(str string) TracableError {
	parent := *e
	parent.root = &allocError{s: str, root: nil, frame: parent.frame, info: parent.info}
	return &parent
}
