package utils

import (
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	Code int
	Data map[string]any
}

type Claims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}
