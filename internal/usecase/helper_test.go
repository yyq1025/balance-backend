package usecase

import (
	"context"
	"errors"
	"reflect"
)

var ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()

var errInternalServErr = errors.New("cannot get networks")
