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
		_ "UAS/docs"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Prestasi Mahasiswa API
// @version 1.0
// @description Backend Sistem Prestasi Mahasiswa (PostgreSQL + MongoDB)
// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization


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
	achievementRepo := repository.NewAchievementRepository(
    database.MongoDB.Collection("achievements"),)
	refRepo := repository.NewAchievementReferenceRepository(db)
	studentRepo := repository.NewStudentRepository(database.DB)
	lecturerRepo := repository.NewLecturerRepository(db)
	reportRepo := repository.NewReportRepository(db)

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

	// Route Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	route.AuthRoute(app, userRepo)
	route.UserRoute(app, userRepo)
	route.AchievementRoute(app, achievementRepo, refRepo, studentRepo)
	route.StudentRoute(app, studentRepo)
	route.LecturerRoute(app, lecturerRepo)
	route.ReportRoute(app, reportRepo, achievementRepo)

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
