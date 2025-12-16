package service

import (
	"UAS/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func UpdateStudentAdvisor(c *fiber.Ctx, repo *repository.StudentRepository) error {
	studentID := c.Params("id")

	var input struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi UUID advisor
	if _, err := uuid.Parse(input.AdvisorID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid advisor_id",
		})
	}

	student, err := repo.UpdateAdvisor(studentID, input.AdvisorID)
	if err != nil {
		if err.Error() == "student not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   student,
	})
}

