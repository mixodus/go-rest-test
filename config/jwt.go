package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("the_secret_key")

type JWTClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
