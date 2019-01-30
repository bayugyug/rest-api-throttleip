package utils

import (
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

const (
	TokenAuthSecret = "/v1/api/S3cr3T/benjerry/icecream/choco"
	TokenAuthExpDay = 365
)

type AppJwtConfig struct {
	TokenAuth *jwtauth.JWTAuth `json:",omitempty"`
}

func NewAppJwtConfig() *AppJwtConfig {
	return &AppJwtConfig{
		TokenAuth: jwtauth.New("HS256", []byte(TokenAuthSecret), nil),
	}
}

func (t *AppJwtConfig) GenToken(claims jwt.MapClaims) (string, error) {
	_, tokenString, err := t.TokenAuth.Encode(claims)
	log.Println("jwt token:", tokenString, err)
	return tokenString, err
}
