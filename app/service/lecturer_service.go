package service

import (
	"UAS/app/repository"

	"github.com/gofiber/fiber/v2"
)

func GetLecturers(c *fiber.Ctx, repo *repository.LecturerRepository) error {
	lecturers, err := repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get lecturers",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": lecturers,
	})
}
