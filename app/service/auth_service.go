package service

import (
	"UAS/app/model"
	"UAS/app/repository"
	"UAS/app/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx, repo *repository.UserRepository) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := repo.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "akun tidak aktif"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	perms, err := repo.GetPermissionsByUserID(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengambil permissions"})
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.RoleID.String(), perms)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal generate token"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"fullName": user.FullName,
				"roleId":   user.RoleID,
			},
		},
	})
}

func RefreshToken(c *fiber.Ctx, repo *repository.UserRepository) error {
	// Ambil refresh token dari body
	var body struct {
		Token string `json:"refreshToken"`
	}
	if err := c.BodyParser(&body); err != nil || body.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token required",
		})
	}

	// Validasi token
	claims, err := utils.ValidateJWT(body.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired refresh token",
		})
	}

	// Ambil user dari DB
	user, err := repo.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Generate token baru
	token, err := utils.GenerateJWT(user.ID.String(), user.RoleID.String(), claims.Permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": token,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Masukkan token ke blacklist
	utils.AddToBlacklist(token, claims.ExpiresAt.Time)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "logout berhasil",
	})
}


