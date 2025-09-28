package jwt

import (
	"ai-software-copyright-server/internal/application/param/request"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWT struct {
	signingKey []byte
}

func NewJWT(signingKey string) *JWT {
	return &JWT{
		[]byte(signingKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.signingKey)
}

func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.signingKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.Wrap(err, "token 解析失败")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.Wrap(err, "登录状态已失效")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.Wrap(err, "token 未激活")
			} else {
				return nil, errors.Wrap(err, "token 无法处理")
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.Wrap(err, "token 无法处理")
	} else {
		return nil, errors.Wrap(err, "token 无法处理")
	}
}
