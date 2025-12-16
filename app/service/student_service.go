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

func GetStudentByID(c *fiber.Ctx, repo *repository.StudentRepository) error {
	studentID := c.Params("id")

	student, err := repo.GetByID(studentID)
	if err != nil {
		if err.Error() == "student not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "student not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": student,
	})
}

