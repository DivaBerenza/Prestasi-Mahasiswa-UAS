package service

import (
	"UAS/app/repository"
	"github.com/gofiber/fiber/v2"
)

func GetStudents(c *fiber.Ctx, repo *repository.StudentRepository) error {
	students, err := repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": students,
	})
}
