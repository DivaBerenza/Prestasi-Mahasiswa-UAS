package route

import (
	"UAS/app/repository"
	"UAS/app/service"
	"UAS/middleware"

	"github.com/gofiber/fiber/v2"
)

// AuthRoute menangani login
func AuthRoute(app *fiber.App, repo *repository.UserRepository) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		return service.Login(c, repo)
	})

	auth.Post("/refresh", middleware.JWTBlacklistMiddleware() ,func(c *fiber.Ctx) error {
		return service.RefreshToken(c, repo)
	})

	auth.Get("/profile", middleware.JWTBlacklistMiddleware(), func(c *fiber.Ctx) error {
		return service.Profile(c, repo)
	})

	auth.Post("/logout", service.Logout)

}

// UserRoute menangani CRUD user (admin only)
func UserRoute(app *fiber.App, repo *repository.UserRepository) {
	users := app.Group("/api/v1/users")

	users.Use(middleware.JWTBlacklistMiddleware())

	// GET all users - hanya admin dengan permission user:read
	users.Get("/", middleware.RBACMiddleware("user:manage"), func(c *fiber.Ctx) error {
		return service.GetAllUsers(c, repo)
	})

	// GET user by ID - admin dan user bisa baca sendiri
	users.Get("/:id", middleware.RBACMiddleware("user:manage"), func(c *fiber.Ctx) error {
		return service.GetUserByID(c, repo)
	})

	// POST create user
	users.Post("/", middleware.RBACMiddleware("user:manage"), func(c *fiber.Ctx) error {
		return service.CreateUser(c, repo)
	})

	// PUT update user
	users.Put("/:id", middleware.RBACMiddleware("user:manage"), func(c *fiber.Ctx) error {
		return service.UpdateUser(c, repo)
	})

	// DELETE user
	users.Delete("/:id", middleware.RBACMiddleware("user:manage"), func(c *fiber.Ctx) error {
		return service.DeleteUser(c, repo)
	})

	// // PUT untuk update password
	users.Put("/:id/role", middleware.RBACMiddleware("user:manage"),func(c *fiber.Ctx) error {
	return service.UpdatePassword(c, repo)
	})
}

func AchievementRoute(app *fiber.App, achievementRepo *repository.AchievementRepository, refRepo *repository.AchievementReferenceRepository, studentRepo *repository.StudentRepository) {
	ach := app.Group("/api/v1/achievements")
	ach.Get("/", func(c *fiber.Ctx) error {
		return service.ListAchievements(c, achievementRepo)
	})

	ach.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAchievementDetail(c, achievementRepo) 
	})

	ach.Post("/", middleware.JWTBlacklistMiddleware(),func(c *fiber.Ctx) error {
    	return service.CreateAchievement(c, achievementRepo, refRepo, studentRepo)
	})

}

func StudentRoute(app *fiber.App, repo *repository.StudentRepository) {
	students := app.Group("/api/v1/students")

	students.Use(middleware.JWTBlacklistMiddleware()) // jika perlu auth

	// GET all students
	students.Get("/", func(c *fiber.Ctx) error {
		return service.ListStudents(c, repo)
	})

	
}

func LecturerRoute(app *fiber.App, repo *repository.LecturerRepository) {
	lec := app.Group("/api/v1/lecturers", middleware.LecturerOnly)

	lec.Get("/", func(c *fiber.Ctx) error {
		return service.GetLecturers(c, repo)
	})

	lec.Get("/:id/advisees", func(c *fiber.Ctx) error {
		return service.GetLecturerAdvisees(c, repo)
	})
}

