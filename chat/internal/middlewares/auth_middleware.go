package middlewares

import (
	"context"
	"fmt"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddlewareInterface interface{}

type AuthMiddlewareStruct struct{}

func NewAuthMiddleware() *AuthMiddlewareStruct {
	return &AuthMiddlewareStruct{}
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func (m *AuthMiddlewareStruct) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := extractBearerToken(c)
		if jwtToken == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "not set jwt token"})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			fmt.Println("jwtSecret:", string([]byte(os.Getenv("JWT_SECRET"))))
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(403, gin.H{"error": "invalid jwt token", "detail": err.Error()})
			return
		}

		idVal := claims["sub"].(float64)
		userID := int(idVal)
		ctx := context.WithValue(c.Request.Context(), jwtinfo_svc.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		email := claims["email"].(string)
		ctx = context.WithValue(c.Request.Context(), jwtinfo_svc.EmailKey, email)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
