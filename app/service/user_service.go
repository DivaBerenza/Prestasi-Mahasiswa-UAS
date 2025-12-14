// service/user_service.go
package service

import (
	"github.com/gofiber/fiber/v2"
	"UAS/app/repository"
)

func GetAllUsers(c *fiber.Ctx, repo *repository.UserRepository) error {
	users, err := repo.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "gagal mengambil data users",
		})
	}

	// Hanya kirim field penting
	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
			"fullName": u.FullName,
			"roleId":   u.RoleID,
			"isActive": u.IsActive,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetUserByID(c *fiber.Ctx, repo *repository.UserRepository) error {
	id := c.Params("id")
	user, err := repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Jika bukan admin, hanya bisa lihat sendiri
	roleID := c.Locals("roleID").(string)
	userID := c.Locals("userID").(string)
	if roleID != "admin-role-uuid" && user.ID.String() != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"fullName": user.FullName,
			"roleId":   user.RoleID,
			"isActive": user.IsActive,
		},
	})
}
