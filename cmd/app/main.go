package main

import (
	"errors"
	"log"
	"os"

	"kpi-backend/internal/database"
	"kpi-backend/internal/registry"
	"kpi-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		// 100 MB request body limit
		BodyLimit: 100 * 1024 * 1024,
		// EnableTrustedProxyCheck: true,
		// Structured JSON error responses for all fiber.NewError / panic recoveries
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var fe *fiber.Error
			if errors.As(err, &fe) {
				code = fe.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://159.65.7.0:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	//health
	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	db := database.ConnectAndMigrate()

	// Initialize Registry and Handlers
	reg := registry.NewRegistry(db)
	handlers := reg.NewAppHandlers()
	routes.RegisterRoutes(app, handlers)

	// ─── Start server ─────────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("Starting KPI backend on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
