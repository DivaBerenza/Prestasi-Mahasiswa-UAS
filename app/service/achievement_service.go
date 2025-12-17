package service

import (
	"net/http"
	"UAS/app/repository"
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
