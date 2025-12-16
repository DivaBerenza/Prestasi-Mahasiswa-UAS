package middleware

import (
	"context"
	"net/http"
	"strings"

	"UAS/app/utils"
	"github.com/gofiber/fiber/v2"
)

type ctxKey string

const (
	UserIDKey ctxKey = "user_id"
	RoleIDKey ctxKey = "role_id"
	PermKey   ctxKey = "permissions"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)
		ctx = context.WithValue(ctx, PermKey, claims.Permissions)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JWTBlacklistMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Set ke context Fiber supaya service bisa ambil
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.RoleID)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

func LecturerOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "lecturer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied: lecturer only",
		})
	}
	return c.Next()
}


