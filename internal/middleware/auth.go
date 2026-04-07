package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"sub"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
	FullName string `json:"name"`
	jwt.RegisteredClaims
}

func Auth() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-change-me"
	}

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || claims.TenantID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid claims")
		}

		c.Locals("userId", claims.UserID)
		c.Locals("tenantId", claims.TenantID)
		c.Locals("role", claims.Role)
		c.Locals("fullName", claims.FullName)

		return c.Next()
	}
}

