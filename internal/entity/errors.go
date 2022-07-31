package entity

import "errors"

var (
	ErrGetNetwork      = errors.New("cannot get networks")
	ErrNetworkNotFound = errors.New("network not found")
	ErrAddWallet       = errors.New("cannot add wallet")
	ErrGetBalance      = errors.New("cannot get balance")
	ErrFindWallet      = errors.New("cannot find wallet")
	ErrDeleteWallet    = errors.New("cannot delete wallet")
)
