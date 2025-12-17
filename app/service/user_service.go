// service/user_service.go
package service

import (
	"github.com/gofiber/fiber/v2"
	
	"UAS/app/repository"
	"UAS/app/model"
	"UAS/app/utils"
	"github.com/google/uuid"
)

func GetAllUsers(c *fiber.Ctx, repo *repository.UserRepository) error {
	users, err := repo.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "gagal mengambil data users",
		})
	}

	// Hanya kirim field penting
	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
			"fullName": u.FullName,
			"roleId":   u.RoleID,
			"isActive": u.IsActive,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetUserByID(c *fiber.Ctx, repo *repository.UserRepository) error {
	id := c.Params("id")

	user, err := repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// User sudah diizinkan oleh middleware, langsung return data
	return c.JSON(fiber.Map{
		"status": "success",
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"fullName": user.FullName,
			"roleId":   user.RoleID,
			"isActive": user.IsActive,
		},
	})
}

func CreateUser(c *fiber.Ctx, repo *repository.UserRepository) error {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
		IsActive bool   `json:"isActive"`

		// lecturer
		LecturerID string `json:"lecturerId"`
		Department string `json:"department"`

		// student
		StudentID    string `json:"studentId"`
		ProgramStudy string `json:"programStudy"`
		AcademicYear string `json:"academicYear"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "hash failed"})
	}

	roleUUID, err := uuid.Parse(input.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid roleId"})
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		FullName: input.FullName,
		RoleID:   roleUUID,
		IsActive: input.IsActive,
	}

	// ===== lecturer ptr =====
	var lecturerIDPtr *string
	var departmentPtr *string

	if input.LecturerID != "" {
		lecturerIDPtr = &input.LecturerID
	}
	if input.Department != "" {
		departmentPtr = &input.Department
	}

	// ===== student ptr =====
	var studentIDPtr *string
	var programStudyPtr *string
	var academicYearPtr *string

	if input.StudentID != "" {
		studentIDPtr = &input.StudentID
	}
	if input.ProgramStudy != "" {
		programStudyPtr = &input.ProgramStudy
	}
	if input.AcademicYear != "" {
		academicYearPtr = &input.AcademicYear
	}

	newUser, err := repo.CreateUser(
		user,
		lecturerIDPtr,
		departmentPtr,
		studentIDPtr,
		programStudyPtr,
		academicYearPtr,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status": "success",
		"data":   newUser,
	})
}



func UpdateUser(c *fiber.Ctx, repo *repository.UserRepository) error {
	// Ambil ID dari URL
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	// Parse body JSON
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
		IsActive bool   `json:"isActive"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Parse RoleID ke uuid.UUID
	roleUUID, err := uuid.Parse(input.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid roleId",
		})
	}

	user := &model.User{
		ID:       userID,
		Username: input.Username,
		Email:    input.Email,
		FullName: input.FullName,
		RoleID:   roleUUID,
		IsActive: input.IsActive,
	}

	updatedUser, err := repo.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": map[string]interface{}{
			"id":       updatedUser.ID,
			"username": updatedUser.Username,
			"email":    updatedUser.Email,
			"fullName": updatedUser.FullName,
			"roleId":   updatedUser.RoleID,
			"isActive": updatedUser.IsActive,
		},
	})
}

func DeleteUser(c *fiber.Ctx, repo *repository.UserRepository) error {
	// Parse ID
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	// Delete user
	err = repo.DeleteUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}


func UpdatePassword(c *fiber.Ctx, repo *repository.UserRepository) error {
	// Parse user ID
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	// Parse body JSON
	var body struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid password",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	// Update password di repository
	updatedUser, err := repo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update password",
		})
	}

	// Response tanpa mengirim password
	return c.JSON(fiber.Map{
		"status": "success",
		"data": map[string]interface{}{
			"id":       updatedUser.ID,
			"username": updatedUser.Username,
			"email":    updatedUser.Email,
			"fullName": updatedUser.FullName,
			"roleId":   updatedUser.RoleID,
			"isActive": updatedUser.IsActive,
		},
	})
}

func Profile(c *fiber.Ctx, repo *repository.UserRepository) error {
    userID, ok := c.Locals("userID").(string)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user"})
    }

    user, err := repo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    return c.JSON(fiber.Map{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "fullName": user.FullName,
        "roleId":   user.RoleID,
        "isActive": user.IsActive,
    })
}







