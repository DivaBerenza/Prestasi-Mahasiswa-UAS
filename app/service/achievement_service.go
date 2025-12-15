package service

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"UAS/app/repository"
)

func ListAchievements(c *fiber.Ctx, repo *repository.AchievementRepository) error {
	// Ambil role dari context (misal di-set oleh auth middleware)
	roleValue := c.Locals("role")
	role, ok := roleValue.(string)
	if !ok {
		role = "student" // default role jika tidak ada di context
	}

	// Ambil query params
	studentId := c.Query("studentId", "")
	status := c.Query("status", "")
	achType := c.Query("type", "")

	// Pagination params
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	achievements, err := repo.GetAchievements(role, studentId, status, achType, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get achievements",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": achievements,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"count": len(achievements),
		},
	})
}

func GetAchievementDetail(c *fiber.Ctx, repo *repository.AchievementRepository) error {
    id := c.Params("id")

    achievement, err := repo.GetAchievementByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Achievement not found",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "data": achievement,
    })
}

