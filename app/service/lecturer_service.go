package service

import (
	"UAS/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func GetLecturerAdvisees(c *fiber.Ctx, repo *repository.LecturerRepository) error {
	lecturerID := c.Params("id")

	// validasi UUID
	if _, err := uuid.Parse(lecturerID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid lecturer id",
		})
	}

	advisees, err := repo.GetAdvisees(lecturerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": advisees,
	})
}


