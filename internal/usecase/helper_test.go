package usecase_test

import (
	"context"
	"errors"
	"reflect"
)

var ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()

var errInternalServErr = errors.New("cannot get networks")

type test struct {
	name string
	mock func()
	res  any
	err  error
}
