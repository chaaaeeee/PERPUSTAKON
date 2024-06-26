package handlers 

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"


	"perpustakaan/middleware"
	"perpustakaan/config"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler() (*Handler, error) {
	db, err := config.Connect()
	if err != nil {
		return nil, err
	}

	return &Handler{
		DB: db,
	}, nil
}

func SetupRoutes(app *fiber.App) {
	handler, err := NewHandler()

	// reminder that im the goat
	app.Use(func (c *fiber.Ctx) error {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"success": false,
					"message": "Error connecting to the database",
					"code": err.Error(),
				},
			)
		}

		return c.Next()
	})

	//swagger 
	app.Get("/swagger/*", swagger.HandlerDefault)

	// books

	// user
	app.Get("/getBooks", handler.GetBooks)
	app.Get("/getBookById/:id", handler.GetBook)
	app.Get("/getBookByTitle/:title", handler.SearchBook)

	// librarian
	app.Post("/addBook", middleware.OnlyLibrarian, handler.AddBook)
	app.Delete("/deleteBook", middleware.OnlyLibrarian, handler.DeleteBook)

	// users

	// user
	app.Post("/signupHandler", middleware.NotLoggedIn, handler.SignupHandler)
	app.Post("/loginHandler", middleware.NotLoggedIn, handler.LoginHandler)
	app.Get("/logoutHandler", handler.LogoutHandler)

	// admin
	app.Get("/getUsers", handler.GetUsers)
	app.Get("/getUserById/:id", handler.GetUser)
	app.Post("/addUser", handler.AddUser)
	app.Delete("/deleteUser", handler.DeleteUser)

	// borrow

	// librarian
	app.Post("/borrowBook", middleware.OnlyLibrarian, handler.BorrowBook)
	app.Post("/returnBook", middleware.OnlyLibrarian, handler.ReturnBook)

	app.Get("/login", handler.Login)

	app.Get("/librarian/dashboard", handler.LibrarianDashboard)
	app.Get("/librarian/bookList", handler.LibrarianBookList)
	app.Get("/librarian/userList", handler.LibrarianUserList)
	app.Get("/librarian/addBook", handler.LibrarianAddBook)
	app.Get("/librarian/deleteBook", handler.LibrarianDeleteBook)
	app.Get("/librarian/borrowBook", handler.LibrarianBorrowBook)
	app.Get("/librarian/returnBook", handler.LibrarianReturnBook)

	app.Get("/admin/dashboard", handler.AdminDashboard)
	app.Get("/admin/userList", handler.AdminUserList)
	app.Get("/admin/addUser", handler.AdminAddUser)
	app.Get("/admin/deleteUser", handler.AdminDeleteUser)

	app.Get("/user/dashboard", handler.UserDashboard)
	app.Get("/user/bookList", handler.UserBookList)

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("./frontend/views/login.html", nil)
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("./frontend/views/signup.html", nil)
	})
}
