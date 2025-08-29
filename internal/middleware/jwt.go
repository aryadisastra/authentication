package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"result": false, "data": gin.H{}, "message": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(h[len("Bearer "):])
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"result": false, "data": gin.H{}, "message": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"result": false, "data": gin.H{}, "message": "invalid claims"})
			return
		}
		c.Set("user_id", claims["sub"])
		c.Set("role_id", claims["role_id"])
		c.Set("role_code", claims["role_code"])
		c.Set("role_name", claims["role_name"])
		c.Next()
	}
}

func RequireRoles(allowed ...string) gin.HandlerFunc {
	allowedSet := map[string]struct{}{}
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}
	return func(c *gin.Context) {
		roleAny, _ := c.Get("role_code")
		role, _ := roleAny.(string)
		if _, ok := allowedSet[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"result": false, "data": gin.H{}, "message": "forbidden"})
			return
		}
		c.Next()
	}
}
