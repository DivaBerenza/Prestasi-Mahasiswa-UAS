package service

import (
	"UAS/app/model"
	"UAS/app/repository"
	"context"

	"github.com/gofiber/fiber/v2"
)

// GetReportStatistics godoc
// @Summary Statistik prestasi keseluruhan
// @Description Admin melihat statistik seluruh prestasi
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.ReportStatistics
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /reports/statistics [get]
func GetReportStatistics(
	reportRepo *repository.ReportRepository,
	achievementRepo *repository.AchievementRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		// 1️⃣ Total mahasiswa (PostgreSQL)
		totalStudents, err := reportRepo.GetStudentStatistics()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to get student statistics",
			})
		}

		// 2️⃣ Statistik prestasi (MongoDB)
		totalAch, submitted, verified, rejected, err :=
			achievementRepo.GetStatistics(ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to get achievement statistics",
			})
		}

		// 3️⃣ Gabungkan hasil
		report := model.ReportStatistics{
			TotalStudents:     totalStudents,
			TotalAchievements: int(totalAch),
			SubmittedCount:    int(submitted),
			VerifiedCount:     int(verified),
			RejectedCount:     int(rejected),
		}

		return c.JSON(fiber.Map{
			"data": report,
		})
	}
}

// GetStudentReport godoc
// @Summary Statistik prestasi mahasiswa
// @Description
// Mahasiswa: hanya data sendiri
// Dosen Wali: data mahasiswa bimbingan
// Admin: semua mahasiswa
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "NIM Mahasiswa"
// @Success 200 {object} model.ReportStudent
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reports/student/{id} [get]
func GetStudentReport(
	reportRepo *repository.ReportRepository,
	achievementRepo *repository.AchievementRepository,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		nim := c.Params("id")

		// 1️⃣ Ambil data dasar mahasiswa (Postgres)
		report, err := reportRepo.GetStudentBase(nim)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		if report == nil {
			return fiber.NewError(404, "student not found")
		}

		// 2️⃣ Hitung prestasi (Mongo pakai NIM)
		total, submitted, verified, rejected, err :=
			achievementRepo.GetStatisticsByStudentID(
				c.Context(),
				nim, // ✅ PAKAI NIM
			)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		// 3️⃣ Mapping hasil
		report.TotalAchievements = int(total)
		report.SubmittedCount = int(submitted)
		report.VerifiedCount = int(verified)
		report.RejectedCount = int(rejected)

		return c.JSON(fiber.Map{
			"data": report,
		})
	}
}




