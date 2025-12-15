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
