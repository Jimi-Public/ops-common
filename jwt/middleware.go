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

	"github.com/gin-gonic/gin"

	"github.com/Jimi-Public/ops-common/response"
)

const AuthHeader = "Authorization"

func JWTMiddleware() gin.HandlerFunc {
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

		c.Set("claims", tokenClaims)
		c.Next()
	}
}
