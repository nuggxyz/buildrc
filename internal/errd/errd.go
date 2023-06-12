package errd

import "errors"

func As[I error](err error) (I, bool) {
	var target I
	ok := errors.As(err, &target)
	return target, ok
}
