/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/10
@IDE:     GoLand
@File:    middleware.go
*/

package jwt

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Jimi-Public/ops-common/response"
)

const AuthHeader = "Authorization"
const ContextClaims = "claims"
const ContextAccountName = "AccountName"

// JWTMiddleware TODO  RBAC 认证。 用户认证和接口鉴权
func JWTMiddleware() gin.HandlerFunc {
	// FIXME 生成token 设置secret 在这里可能拿不到, 故使用默认值
	j := NewToken()
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(AuthHeader)
		if authHeader == "" {
			response.Resp(c, response.AuthFail, "Authentication failure", response.Data{})
			c.Abort()
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			response.Resp(c, response.AuthFail, "Authorization header is malformed", response.Data{})
			c.Abort()
			return
		}

		token := authHeaderParts[1]
		tokenClaims, err := j.ParseToken(token)
		if err != nil {
			response.Resp(c, response.AuthFail, "Token is invalid", response.Data{})
			c.Abort()
			return
		}
		if tokenClaims.ExpiresAt < time.Now().Unix() {
			response.Resp(c, response.TokenExpire, "Token is Expire", response.Data{})
			c.Abort()
		}
		// Context 插入 上下文
		c.Set(ContextClaims, tokenClaims)
		c.Set(ContextAccountName, tokenClaims.AccountName)
		c.Set(AuthHeader, authHeader)
		c.Next()
	}
}
