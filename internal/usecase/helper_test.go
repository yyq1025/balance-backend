package usecase_test

import (
	"errors"
)

var errInternalServErr = errors.New("cannot get networks")

type test struct {
	name string
	mock func()
	res  any
	err  error
}
