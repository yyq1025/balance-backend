package usecase_test

import (
	"errors"
)

var errInternalServerErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  any
	err  error
}
