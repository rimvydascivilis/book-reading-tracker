package rest_test

import (
	"github.com/golang-jwt/jwt/v5"
)

var mockJWTToken = jwt.Token{
	Claims: jwt.MapClaims{
		"id": float64(1),
	},
}
