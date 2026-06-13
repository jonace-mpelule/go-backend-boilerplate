package main

import (
	"log"

	"github.com/username/project-name/internal/app"
)

// @title Go Backend Boilerplate API
// @version 1.0
// @description Production-ready modular monolith boilerplate built with Chi, Ent, Postgres, Redis, and structured observability.
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide a bearer token using the format: Bearer <token>
func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal("error starting app:", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
