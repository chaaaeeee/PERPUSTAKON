package repository

import (
	"github.com/gofiber/fiber/v2"
	"perpustakaan/middleware"
)

func (repo *Repository) SetupRoutes(app *fiber.App) {
	// admin
	app.Get("/users", middleware.OnlyAdmin, repo.GetUsers)
	// app.Get("/users", middleware.OnlyLibrarian(GetUsers))
	// app.Get("/users", middleware.OnlyAdmin(GetUsers))

	// librarian

	// user
	app.Post("/signup", repo.Signup)
	app.Post("/login", repo.Login)
}

