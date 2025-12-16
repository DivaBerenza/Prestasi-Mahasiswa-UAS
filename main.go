package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"UAS/config"
	"UAS/database"
	"UAS/app/repository"
	"UAS/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect ke DB
	db := database.Connect(cfg.DB_DSN)
	fmt.Println("✅ Database connected successfully!")

	database.ConnectMongo()
	fmt.Println("✅ MongoDB connected successfully!")

	// =========================

	// Init repository
	userRepo := repository.NewUserRepository(db)
	achievementRepo := repository.NewAchievementRepository(database.MongoDB)
	refRepo := repository.NewAchievementReferenceRepository(db)
	studentRepo := repository.NewStudentRepository(database.DB)
	lecturerRepo := repository.NewLecturerRepository(db)



	// Init Fiber
	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 API Running. Database Connected Successfully.")
	})

	// Auth route login
	route.AuthRoute(app, userRepo)

	// User route (CRUD users, admin only)
	route.UserRoute(app, userRepo)

	route.AchievementRoute(app, achievementRepo, refRepo, studentRepo)

	route.StudentRoute(app, studentRepo)
	route.LecturerRoute(app, lecturerRepo)

	// Channel untuk Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Jalankan server di goroutine
	go func() {
		port := cfg.AppPort
		log.Printf("🚀 Server is running on http://localhost:%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Tunggu Ctrl+C
	<-c
	fmt.Println("\n🛑 Server stopped gracefully.")
	os.Exit(0)
}
