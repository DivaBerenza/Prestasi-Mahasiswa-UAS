package service

import (
	"net/http"
	"time"
	"fmt"
	"os"
	"strings"

	"UAS/app/model"
	"path/filepath"
	"UAS/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


// ============================
// LIST ACHIEVEMENTS
// ============================
func ListAchievements(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	studentRepo *repository.StudentRepository,
) error {

	userIDAny := c.Locals("userID")
	permsAny := c.Locals("permissions")

	if userIDAny == nil || permsAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID := userIDAny.(string)
	perms := permsAny.([]string)

	isStudent := false
	for _, p := range perms {
		if p == "achievement:create" {
			isStudent = true
			break
		}
	}

	// ===== MAHASISWA =====
	if isStudent {
		students, err := studentRepo.GetByUserID(userID)
		if err != nil || len(students) == 0 {
			return c.JSON(fiber.Map{"data": []interface{}{}})
		}

		var studentIDs []string
		for _, s := range students {
			studentIDs = append(studentIDs, s.StudentID)
		}

		achievements, err := achievementRepo.GetByStudentID(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"data": achievements})
	}

	// ===== DOSEN / ADMIN =====
	achievements, err := achievementRepo.GetVerifiedOnly()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": achievements})
}


// ============================
// DETAIL
// ============================
func GetAchievementDetail(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	achievement, err := achievementRepo.GetByID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if achievement == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "achievement not found"})
	}

	return c.JSON(fiber.Map{
		"status":      "success",
		"achievement": achievement,
	})
}


// ============================
// CREATE
// ============================
func CreateAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	refRepo *repository.AchievementReferenceRepository,
	studentRepo *repository.StudentRepository,
) error {

	userID := c.Locals("userID").(string)

	students, err := studentRepo.GetByUserID(userID)
	if err != nil || len(students) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "student profile not found"})
	}

	student := students[0]

	var input model.Achievement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	input.StudentID = student.StudentID
	input.Status = string(model.StatusDraft)
	input.Points = 0
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	created, err := achievementRepo.Create(&input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	mongoID := created.ID.Hex()

	ref, err := refRepo.CreateReference(student.ID, mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":      "success",
		"achievement": created,
		"reference":   ref,
	})
}


// ============================
// UPDATE (DRAFT ONLY)
// ============================
func UpdateAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	studentRepo *repository.StudentRepository,
) error {

	userID := c.Locals("userID").(string)
	id := c.Params("id")

	students, err := studentRepo.GetByUserID(userID)
	if err != nil || len(students) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "student profile not found"})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	achievement, err := achievementRepo.GetByID(id)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if achievement.StudentID != students[0].StudentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if achievement.Status != string(model.StatusDraft) {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	updated, err := achievementRepo.Update(id, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "achievement": updated})
}


// ============================
// DELETE
// ============================
func DeleteAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	studentRepo *repository.StudentRepository,
) error {

	userID := c.Locals("userID").(string)
	id := c.Params("id")

	students, err := studentRepo.GetByUserID(userID)
	if err != nil || len(students) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "student profile not found"})
	}

	achievement, err := achievementRepo.GetByID(id)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if achievement.StudentID != students[0].StudentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if achievement.Status != string(model.StatusDraft) {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if err := achievementRepo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}


// ============================
// SUBMIT (MAHASISWA)
// ============================
func SubmitAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	refRepo *repository.AchievementReferenceRepository,
	studentRepo *repository.StudentRepository,
) error {

	userID := c.Locals("userID").(string)
	id := c.Params("id")

	students, err := studentRepo.GetByUserID(userID)
	if err != nil || len(students) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "student profile not found"})
	}

	achievement, err := achievementRepo.GetByID(id)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if achievement.StudentID != students[0].StudentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if achievement.Status != string(model.StatusDraft) {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	updated, err := achievementRepo.UpdateStatus(id, string(model.StatusSubmitted))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := refRepo.Submit(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "achievement": updated})
}


// ============================
// VERIFY (DOSEN)
// ============================
func VerifyAchievement(
	c *fiber.Ctx,
	refRepo *repository.AchievementReferenceRepository,
) error {

	mongoID := c.Params("id")
	lecturerID := c.Locals("userID").(string)

	uid, err := uuid.Parse(lecturerID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	if err := refRepo.Verify(mongoID, uid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement verified"})
}


// ============================
// REJECT (DOSEN)
// ============================
func RejectAchievement(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
	refRepo *repository.AchievementReferenceRepository,
) error {

	id := c.Params("id")

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	achievement, err := achievementRepo.GetByID(id)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if achievement.Status != string(model.StatusSubmitted) {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
	}

	updated, err := achievementRepo.UpdateStatus(id, string(model.StatusRejected))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := refRepo.Reject(id, body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "achievement": updated})
}


// ============================
// HISTORY
// ============================
func GetAchievementHistory(
	c *fiber.Ctx,
	refRepo *repository.AchievementReferenceRepository,
) error {

	id := c.Params("id")

	history, err := refRepo.StatusHistory(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if history == nil {
		return c.Status(404).JSON(fiber.Map{"error": "history not found"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   []model.AchievementStatusHistory{*history},
	})
}

// service/achievement_service.go
func UploadAchievementAttachment(
	c *fiber.Ctx,
	achievementRepo *repository.AchievementRepository,
) error {

	achievementID := c.Params("id")

	// Ambil file
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	// Validasi ekstensi
	allowedExt := map[string]bool{
		".png": true,
		".jpg": true,
		".jpeg": true,
		".pdf": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return fiber.NewError(fiber.StatusBadRequest, "invalid file type")
	}

	// Folder upload
	uploadDir := fmt.Sprintf("./uploads/achievements/%s", achievementID)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Nama file unik
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Simpan file
	if err := c.SaveFile(file, filePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	attachment := model.Attachment{
		FileName: file.Filename,
		FilePath: filePath,
		FileType: ext,
		UploadedAt: time.Now(),
	}

	// Push attachment ke achievement
	if err := achievementRepo.AddAttachment(achievementID, attachment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "attachment uploaded successfully",
		"data": attachment,
	})
}
