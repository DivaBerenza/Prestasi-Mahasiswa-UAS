package service

import (
	"UAS/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetLecturers godoc
// @Summary Get all lecturers
// @Description Mengambil daftar seluruh dosen
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturers [get]
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

// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Mengambil daftar mahasiswa bimbingan dosen
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lecturer ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturers/{id}/advisees [get]
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


