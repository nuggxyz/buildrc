package terrors

import "github.com/go-faster/errors"

type Framer interface {
	Root() error
	Frame() errors.Frame
	Info() []any
}

// Cause returns first recorded Frame.
func Cause(err error) (oerr error, f errors.Frame, i []any, r bool) {
	for {
		we, ok := err.(Framer)
		if !ok {
			return err, f, i, r
		}
		f = we.Frame()
		i = we.Info()
		r = r || ok
		oerr = err

		err = we.Root()
		if err == nil {
			return oerr, f, i, r
		}
	}
}

func Cause2(err error) (f Framer, r bool) {
	for {
		we, ok := err.(Framer)
		if !ok {
			return
		}
		r = r || ok
		f = we

		err = we.Root()
		if err == nil {
			return
		}
	}
}
