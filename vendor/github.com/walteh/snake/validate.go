package snake

import (
	"reflect"

	"github.com/go-faster/errors"
)

var (
	ErrInvalidMethodSignature = errors.New("invalid method signatured")
)

func commandResponseValidationStrategy(out []reflect.Type) error {

	if len(out) != 1 {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return signature, expected 1, got %d", len(out))
	}

	if !out[0].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q, expected %q", out[0].String(), reflect.TypeOf((*error)(nil)).Elem().String())
	}

	return nil
}

func commandResponseHandleStrategy(out []reflect.Value) ([]*reflect.Value, error) {

	resp := []*reflect.Value{end_of_chain_ptr}

	if !out[0].IsNil() {
		return resp, out[0].Interface().(error)
	}

	return resp, nil
}

func handleArgumentResponse(out []reflect.Value, inter ...any) ([]*reflect.Value, error) {

	res := make([]*reflect.Value, len(inter))

	if !out[len(out)-1].IsNil() {
		// need to fix this TODO
		return nil, out[len(out)-1].Interface().(error)
	}

	for i, v := range inter {
		if out[i].Type() != reflect.TypeOf(v).Elem() {
			return nil, errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q, expected %q", out[i].String(), reflect.TypeOf(v).Elem().String())
		}
		res[i] = &out[i]
	}

	return res, nil
}

func validateArgumentResponse(out []reflect.Type, inter ...any) error {

	if len(out) != len(inter)+1 {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return signature, expected 2, got %d", len(out))
	}

	for i, v := range inter {
		if !out[i].Implements(reflect.TypeOf(v).Elem()) {
			return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q, expected %q", out[i].String(), reflect.TypeOf(v).Elem().String())
		}
	}

	if !out[len(out)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q, expected %q", out[len(out)-1].String(), reflect.TypeOf((*error)(nil)).Elem().String())
	}

	return nil
}

func handle1ArgumentResponse[A any](out []reflect.Value) ([]*reflect.Value, error) {
	return handleArgumentResponse(out, (*A)(nil))
}

func validate1ArgumentResponse[A any](out []reflect.Type) error {
	return validateArgumentResponse(out, (*A)(nil))
}

func validate2ArgumentResponse[A any, B any](out []reflect.Type) error {
	return validateArgumentResponse(out, (*A)(nil), (*B)(nil))
}

func handle2ArgumentResponse[A any, B any](out []reflect.Value) ([]*reflect.Value, error) {
	return handleArgumentResponse(out, (*A)(nil), (*B)(nil))
}

func validate3ArgumentResponse[A any, B any, C any](out []reflect.Type) error {
	return validateArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil))
}

func handle3ArgumentResponse[A any, B any, C any](out []reflect.Value) ([]*reflect.Value, error) {
	return handleArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil))
}

func validate4ArgumentResponse[A any, B any, C any, D any](out []reflect.Type) error {
	return validateArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil), (*D)(nil))
}

func handle4ArgumentResponse[A any, B any, C any, D any](out []reflect.Value) ([]*reflect.Value, error) {
	return handleArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil), (*D)(nil))
}

func validate5ArgumentResponse[A any, B any, C any, D any, E any](out []reflect.Type) error {
	return validateArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil), (*D)(nil), (*E)(nil))
}

func handle5ArgumentResponse[A any, B any, C any, D any, E any](out []reflect.Value) ([]*reflect.Value, error) {
	return handleArgumentResponse(out, (*A)(nil), (*B)(nil), (*C)(nil), (*D)(nil), (*E)(nil))
}
