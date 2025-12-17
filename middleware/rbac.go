package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"UAS/app/utils"
)

// RBACMiddleware menerima permission yang dibutuhkan
func RBACMiddleware(requiredPerms ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing or invalid Authorization header",
            })
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := utils.ValidateJWT(tokenStr)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }

        c.Locals("userID", claims.UserID)
        c.Locals("roleID", claims.RoleID)
        c.Locals("permissions", claims.Permissions)

        hasPerm := false
        for _, p := range claims.Permissions {
            for _, rp := range requiredPerms {
                if p == rp {
                    hasPerm = true
                    break
                }
            }
            if hasPerm {
                break
            }
        }

        if !hasPerm {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }

        return c.Next()
    }
}

