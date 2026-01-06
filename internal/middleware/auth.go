package middleware

import (
	"net/http"
	"strings"

	"github.com/algamoney/api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token não fornecido"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "claims inválidos"})
			return
		}

		c.Set("user_email", claims["sub"])
		c.Set("user_name", claims["nome"])

		if authorities, ok := claims["authorities"].([]interface{}); ok {
			var perms []string
			for _, auth := range authorities {
				if perm, ok := auth.(string); ok {
					perms = append(perms, perm)
				}
			}
			c.Set("authorities", perms)
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities, exists := c.Get("authorities")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
			return
		}

		perms, ok := authorities.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
			return
		}

		hasPermission := false
		for _, p := range perms {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado - permissão necessária: " + permission})
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities, exists := c.Get("authorities")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
			return
		}

		perms, ok := authorities.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
			return
		}

		hasPermission := false
		for _, required := range permissions {
			for _, p := range perms {
				if p == required {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
			return
		}

		c.Next()
	}
}
