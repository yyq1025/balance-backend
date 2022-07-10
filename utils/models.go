package utils

import (
	"context"
)

type Response struct {
	Code int
	Data map[string]any
}

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}
