package middleware

import (
	"context"
	"errors"
	"net/http"
	"photo-web/api"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenStatusChecker interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

// AuthRequired 认证中间件：验证 JWT Token
func AuthRequired(tokenChecker TokenStatusChecker, jwtManager *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 Bearer Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.Error(c, http.StatusUnauthorized, api.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			api.Error(c, http.StatusUnauthorized, api.CodeUnauthorized, "Token 格式错误")
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 检查黑名单
		if tokenChecker.IsTokenBlacklisted(c.Request.Context(), tokenString) {
			api.Error(c, http.StatusUnauthorized, api.CodeTokenExpired, "Token 已失效")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				api.Error(c, http.StatusUnauthorized, api.CodeTokenExpired, "Token 已过期")
			} else {
				api.Error(c, http.StatusUnauthorized, api.CodeUnauthorized, "Token 无效")
			}
			c.Abort()
			return
		}

		// 将用户信息存入 context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
