package service

import (
	"UAS/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetStudents godoc
// @Summary Get all students
// @Description Mengambil daftar seluruh mahasiswa
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students [get]
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

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Mengambil data mahasiswa berdasarkan ID
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/{id} [get]
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

// UpdateStudentAdvisor godoc
// @Summary Update student advisor
// @Description Mengubah dosen wali mahasiswa
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID (UUID)"
// @Param body body object true "Advisor payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/{id}/advisor [put]
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

