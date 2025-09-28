package utils

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

func GenerateToken(tokenId string, tokenType string, userId int64, userType string, expiration time.Duration) (string, int64, error) {
	claims := request.CustomClaims{
		Type:     tokenType,
		UserId:   userId,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenId,
			NotBefore: time.Now().Unix(),                 // 签名生效时间
			ExpiresAt: time.Now().Add(expiration).Unix(), // 过期时间 7天
			Issuer:    "Tool",                            // 签名的发行者
		},
	}
	tokenJwt, err := global.JWT.CreateToken(claims)
	if err != nil {
		return "", 0, errors.Wrap(err, "获取token失败")
	}
	err = global.CACHE.SetCache(fmt.Sprintf("%s_%s_%d_%s", tokenType, userType, userId, tokenId), tokenJwt, expiration)
	if err != nil {
		return "", 0, errors.Wrap(err, "存储token失败")
	}
	return tokenJwt, claims.ExpiresAt, nil
}

func AuthToken(userId int64, userType string) (*common.Token, error) {
	tokenId := uuid.New().String()
	tokenJwt, expiresAt, err := GenerateToken(tokenId, global.AuthToken, userId, userType, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	refreshTokenJwt, _, err := GenerateToken(tokenId, global.RefreshToken, userId, userType, 15*24*time.Hour)
	if err != nil {
		return nil, err
	}
	token := &common.Token{
		AccessToken:  tokenJwt,
		ExpiredIn:    expiresAt,
		RefreshToken: refreshTokenJwt,
	}
	return token, nil
}

func RefreshToken(param request.RefreshTokenParam) (token common.Token, err error) {
	// parseToken 解析token包含的信息
	claims, err := global.JWT.ParseToken(param.Token)
	if err != nil || claims.Type != global.RefreshToken {
		return token, errors.Wrap(err, "Token 已失效")
	}
	checkKey := fmt.Sprintf("%s_%s_%d_%s", global.RefreshToken, claims.UserType, claims.UserId, claims.Id)
	if _, exist := global.CACHE.GetCache(checkKey); exist {
		tokenId := uuid.New().String()
		tokenJwt, expiresAt, err := GenerateToken(tokenId, global.AuthToken, claims.UserId, claims.UserType, 7*24*time.Hour)
		if err != nil {
			return token, err
		}
		refreshTokenJwt, _, err := GenerateToken(tokenId, global.RefreshToken, claims.UserId, claims.UserType, 15*24*time.Hour)
		if err != nil {
			return token, err
		}
		token = common.Token{
			AccessToken:  tokenJwt,
			ExpiredIn:    expiresAt,
			RefreshToken: refreshTokenJwt,
		}
		global.CACHE.DeleteCache(checkKey)
		global.CACHE.DeleteCache(fmt.Sprintf("%s_%s_%d_%s", global.AuthToken, claims.UserType, claims.UserId, claims.Id))
		return token, err
	}
	return token, errors.New("Token 已失效")
}

func RemoveToken(claims *request.CustomClaims) {
	global.CACHE.DeleteCache(fmt.Sprintf("%s_%s_%d_%s", global.AuthToken, claims.UserType, claims.UserId, claims.Id))
	global.CACHE.DeleteCache(fmt.Sprintf("%s_%s_%d_%s", global.RefreshToken, claims.UserType, claims.UserId, claims.Id))
}
