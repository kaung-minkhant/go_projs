package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kaung-minkhant/go_projs/go_fiber_crm/database"
	"github.com/kaung-minkhant/go_projs/go_fiber_crm/lead"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/leads", lead.GetLeads)
	app.Get("/api/v1/leads/:id", lead.GetLead)
	app.Post("/api/v1/leads", lead.NewLead)
	app.Delete("/api/v1/leads/:id", lead.DeleteLead)
}

func main() {
	app := fiber.New()
	database.InitDatabase()

	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated")

	setupRoutes(app)

	app.Listen(":8080")
}
