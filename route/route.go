package route

import (
	"UAS/app/repository"
	"UAS/app/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App, repo *repository.UserRepository) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		return service.Login(c, repo)
	})
}
