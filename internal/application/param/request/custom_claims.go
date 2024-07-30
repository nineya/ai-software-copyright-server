package request

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserId int64  `json:"userId"  form:"userId"`
	Type   string `json:"type" form:"type"`
}
