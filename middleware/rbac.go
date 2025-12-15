package middleware

import (
	"strings"

	"UAS/app/utils"

	"github.com/gofiber/fiber/v2"
)

// RBACMiddleware menerima permission yang dibutuhkan
func RBACMiddleware(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		// Ambil token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Simpan data user di Locals Fiber agar bisa diakses di service
		c.Locals("userID", claims.UserID)
		c.Locals("roleID", claims.RoleID)
		c.Locals("permissions", claims.Permissions)

		// Cek permission
		hasPerm := false
		for _, p := range claims.Permissions {
			if p == requiredPerm {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			// Jika permission yang dibutuhkan adalah user:manage, berarti bukan admin → tolak
			if requiredPerm == "user:manage" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Anda tidak memiliki akses",
				})
			}
			// Default error untuk permission lain
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		// Lanjut ke handler berikutnya
		return c.Next()
	}
}
