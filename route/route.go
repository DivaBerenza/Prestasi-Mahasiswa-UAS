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
	users := app.Group("/api/v1/users", middleware.JWTBlacklistMiddleware())

	// users.Use(middleware.JWTBlacklistMiddleware())

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
	ach := app.Group("/api/v1/achievements", middleware.JWTBlacklistMiddleware())

	ach.Get("/", middleware.RBACMiddleware("user:manage", "achievement:read"),func(c *fiber.Ctx) error {
		return service.ListAchievements(c, achievementRepo, studentRepo)
	})

	ach.Get("/:id", middleware.RBACMiddleware("achievement:read", "user:manage"),func(c *fiber.Ctx) error {
		return service.GetAchievementDetail(c, achievementRepo) 
	})

	ach.Post("/", middleware.RBACMiddleware("achievement:create"), func(c *fiber.Ctx) error {
		return service.CreateAchievement(c, achievementRepo, refRepo, studentRepo)
	})

	ach.Put("/:id", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		return service.UpdateAchievement(c, achievementRepo, studentRepo)
	})

	ach.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), func(c *fiber.Ctx) error {
		return service.DeleteAchievement(c, achievementRepo, studentRepo)
	})

	// Submit for verification (Mahasiswa)
	ach.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), func(c *fiber.Ctx) error {
		return service.SubmitAchievement(c, achievementRepo, refRepo, studentRepo)
	})

	// Verify & Reject (Dosen Wali)
	ach.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		return service.VerifyAchievement(c, refRepo)
	})

	ach.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		return service.RejectAchievement(c, achievementRepo, refRepo)
	}) 

	// History & Attachments
	ach.Get("/:id/history", middleware.RBACMiddleware("achievement:read", "user:manage"), func(c *fiber.Ctx) error {
		return service.GetAchievementHistory(c, refRepo) // function service belum dibuat
	})

	ach.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		return service.UploadAchievementAttachment(c, achievementRepo) // function service belum dibuat
	})

}

func StudentRoute(app *fiber.App, repo *repository.StudentRepository) {
	students := app.Group("/api/v1/students", middleware.JWTBlacklistMiddleware())

	// GET all students
	students.Get("/", func(c *fiber.Ctx) error {
		return service.GetStudents(c, repo)
	})

	students.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetStudentByID(c, repo)
	})

	students.Put("/:id/advisor", middleware.RBACMiddleware("user:manage"),func(c *fiber.Ctx) error {
		return service.UpdateStudentAdvisor(c, repo)
	})

}

func LecturerRoute(app *fiber.App, repo *repository.LecturerRepository) {
	lec := app.Group("/api/v1/lecturers", middleware.JWTBlacklistMiddleware())

	lec.Get("/", func(c *fiber.Ctx) error {
		return service.GetLecturers(c, repo)
	})

	lec.Get("/:id/advisees", middleware.RBACMiddleware("achievement:read", "user:manage"),func(c *fiber.Ctx) error {
		return service.GetLecturerAdvisees(c, repo)
	})
}

func ReportRoute( app *fiber.App, reportRepo *repository.ReportRepository, achievementRepo *repository.AchievementRepository,) {
	r := app.Group( "/api/v1/reports", middleware.JWTBlacklistMiddleware())

	r.Get("/statistics",
		service.GetReportStatistics(reportRepo, achievementRepo),
	)

	r.Get("/student/:id",
		service.GetStudentReport(reportRepo, achievementRepo),
	)
}

