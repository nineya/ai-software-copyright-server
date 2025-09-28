package request

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	Type     string `json:"type" form:"type" label:"Token类型"`      // token 类型
	UserId   int64  `json:"userId"  form:"userId" label:"用户ID"`    // 用户id
	UserType string `json:"userType" form:"userType" label:"用户类型"` // 用户类型
}
