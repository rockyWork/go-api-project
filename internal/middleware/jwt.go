package middleware

import (
	"errors"
	"strings"

	"go-api-project/internal/model"
	"go-api-project/pkg/jwt"
	"go-api-project/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

func JWTAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrExpiredToken) {
				response.Error(c, 401002, "token expired")
			} else {
				response.Unauthorized(c, "invalid token")
			}
			c.Abort()
			return
		}

		if claims.Type != "access" {
			response.Unauthorized(c, "invalid token type")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRoleKey)
		if !exists || role.(int) != model.UserRoleAdmin {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}
	return userID.(uint)
}

func GetRole(c *gin.Context) int {
	role, exists := c.Get(ContextRoleKey)
	if !exists {
		return 0
	}
	return role.(int)
}

func IsAdmin(c *gin.Context) bool {
	return GetRole(c) == model.UserRoleAdmin
}
