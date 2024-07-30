package admin

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/param/request"
	"tool-server/internal/application/param/response"
	"tool-server/internal/application/service"
	"tool-server/internal/global"
	"tool-server/internal/utils"
)

type AuthService struct {
	service.BaseService
}

var onceAuth = sync.Once{}
var authService *AuthService

// 获取单例
func GetAuthService() *AuthService {
	onceAuth.Do(func() {
		authService = new(AuthService)
		authService.Db = global.DB
	})
	return authService
}

func (s *AuthService) Login(param request.AdminLoginParam) (token response.TokenResponse, err error) {
	var admin table.Admin
	exist, err := s.Db.Where("username = ?", param.Username).Get(&admin)
	if err != nil {
		return token, err
	}
	if !exist {
		return token, errors.New("用户不存在或已被禁用")
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), utils.Md5ByBytes(param.Password))
	if err != nil {
		return token, err
	}
	tokenId := uuid.New().String()
	tokenJwt, expiresAt, err := s.generateToken(admin.Id, tokenId, global.AuthToken, 7*24*time.Hour)
	if err != nil {
		return token, err
	}
	refreshTokenJwt, _, err := s.generateToken(admin.Id, tokenId, global.RefreshToken, 15*24*time.Hour)
	if err != nil {
		return token, err
	}
	token = response.TokenResponse{
		AccessToken:  tokenJwt,
		ExpiredIn:    expiresAt,
		RefreshToken: refreshTokenJwt,
	}
	return token, err
}

func (s *AuthService) Logout(claims *request.CustomClaims) {
	global.CACHE.DeleteCache(fmt.Sprintf("%s_%d_%s", global.AuthToken, claims.UserId, claims.Id))
	global.CACHE.DeleteCache(fmt.Sprintf("%s_%d_%s", global.RefreshToken, claims.UserId, claims.Id))
}

func (s *AuthService) RefreshToken(param request.RefreshTokenParam) (token response.TokenResponse, err error) {
	// parseToken 解析token包含的信息
	claims, err := global.JWT.ParseToken(param.Token)
	if err != nil || claims.Type != global.RefreshToken {
		return token, errors.Wrap(err, "Token 已失效")
	}
	checkKey := fmt.Sprintf("%s_%d_%s", global.RefreshToken, claims.UserId, claims.Id)
	if _, exist := global.CACHE.GetCache(checkKey); exist {
		tokenId := uuid.New().String()
		tokenJwt, expiresAt, err := s.generateToken(claims.UserId, tokenId, global.AuthToken, 7*24*time.Hour)
		if err != nil {
			return token, err
		}
		refreshTokenJwt, _, err := s.generateToken(claims.UserId, tokenId, global.RefreshToken, 15*24*time.Hour)
		if err != nil {
			return token, err
		}
		token = response.TokenResponse{
			AccessToken:  tokenJwt,
			ExpiredIn:    expiresAt,
			RefreshToken: refreshTokenJwt,
		}
		global.CACHE.DeleteCache(checkKey)
		global.CACHE.DeleteCache(fmt.Sprintf("%s_%d_%s", global.AuthToken, claims.UserId, claims.Id))
		return token, err
	}
	return token, errors.New("Token 已失效")
}

func (s *AuthService) generateToken(userId int64, tokenId string, tokenType string, expiration time.Duration) (string, int64, error) {
	claims := request.CustomClaims{
		UserId: userId,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenId,
			NotBefore: time.Now().Unix(),                 // 签名生效时间
			ExpiresAt: time.Now().Add(expiration).Unix(), // 过期时间 7天
			Issuer:    "Aurora",                          // 签名的发行者
		},
	}
	tokenJwt, err := global.JWT.CreateToken(claims)
	if err != nil {
		return "", 0, errors.Wrap(err, "获取token失败")
	}
	err = global.CACHE.SetCache(fmt.Sprintf("%s_%d_%s", tokenType, userId, tokenId), tokenJwt, expiration)
	if err != nil {
		return "", 0, errors.Wrap(err, "存储token失败")
	}
	return tokenJwt, claims.ExpiresAt, nil
}
