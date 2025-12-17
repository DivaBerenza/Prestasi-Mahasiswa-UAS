package service

import (
	"net/http"
	"time"

	"UAS/app/repository"
	"UAS/app/model"

	"github.com/gofiber/fiber/v2"
)

// ListAchievements menampilkan daftar prestasi, filtered by role
func ListAchievements(c *fiber.Ctx, achievementRepo *repository.AchievementRepository, studentRepo *repository.StudentRepository) error {
    userID := c.Locals("userID").(string)

    // Ambil mahasiswa yang terkait dengan userID
    students, err := studentRepo.GetByUserID(userID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    if len(students) == 0 {
        return c.JSON(fiber.Map{"data": []interface{}{}})
    }

    var studentIDs []string
    for _, s := range students {
        studentIDs = append(studentIDs, s.StudentID)
    }

    achievements, err := achievementRepo.GetByStudentID(studentIDs)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"data": achievements})
}

func CreateAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	refRepo *repository.AchievementReferenceRepository,
	studentRepo *repository.StudentRepository,
) error {
	userID := c.Locals("userID").(string)

	// Ambil data mahasiswa dari userID
	students, err := studentRepo.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if len(students) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student profile not found"})
	}

	student := students[0]
	studentUUID := student.ID       // UUID untuk PostgreSQL
	studentNIM := student.StudentID // string/NIM untuk MongoDB

	// Parse body request
	var input model.Achievement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set field default
	input.StudentID = studentNIM
	input.Status = string(model.StatusDraft)
	input.Points = 0
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	// Simpan ke MongoDB
	createdAchievement, err := achievementRepo.Create(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil mongoID
	mongoID := createdAchievement.ID.Hex()

	// Simpan reference di PostgreSQL
	ref, err := refRepo.CreateReference(studentUUID, mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed create reference: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "success",
		"achievement": createdAchievement,
		"reference":   ref,
	})
}

