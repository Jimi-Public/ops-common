/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/9
@IDE:     GoLand
@File:    jwt.go
*/

package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenInterface interface {
	GenerateToken(id int, accountName string) (string, error)
	ParseToken(token string) (*Claims, error)
}

type Token struct {
	JwtSecret  string // 加密秘钥
	ExpireTime int    // 多少小时过期
}

func NewToken(f ...func(token *Token)) TokenInterface {
	t := &Token{}
	for _, i := range f {
		i(t)
	}
	// 未赋值则初始化
	if t.ExpireTime == 0 {
		t.ExpireTime = 24
	}
	if t.JwtSecret == "" {
		t.JwtSecret = "secret"
	}
	return t
}

// OptionWithJwtSecret 设置加密salt
func OptionWithJwtSecret(s string) func(token *Token) {
	return func(token *Token) {
		token.JwtSecret = s
	}
}

// OptionWithExpireTime 设置Token过期时间(单位: 小时)
func OptionWithExpireTime(t int) func(token *Token) {
	return func(token *Token) {
		token.ExpireTime = t
	}
}

type Claims struct {
	UserID      int    `json:"user_id"`
	AccountName string `json:"account_name"`
	jwt.StandardClaims
}

// GenerateToken 生成Token
func (t *Token) GenerateToken(id int, accountName string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(t.ExpireTime) * time.Hour)

	claims := Claims{
		UserID:      id,
		AccountName: accountName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "jimi-ops",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(t.JwtSecret))

	return token, err
}

// ParseToken 解析Token
func (t *Token) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.JwtSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
