package repository

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"perpustakaan/models"
	"perpustakaan/service"
)

func (r *Repository) Signup(c *fiber.Ctx) error {
	db := r.DB

	type userStruct struct {
		Username string
		Password string
	}

	var user userStruct
	var dbUser models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err := db.QueryRow("SELECT username FROM users WHERE username=?", user.Username).Scan(&dbUser.Username)
	if err != sql.ErrNoRows { // if username has already taken
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Username is already taken"})
	}

	/*
		if err != nil {
			fmt.Println("atas kah")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	*/

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, hashedPassword, 1)
	if err != nil {
		fmt.Println("sini kah")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User created successfully"})
}

func (r *Repository) Login(c *fiber.Ctx) error {
	db := r.DB

	type userStruct struct {
		Username string
		Password string
	}

	var user userStruct
	var dbUser models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username=?", user.Username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password, &dbUser.Role)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Username is not registered"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Incorrect password"})
	}
	// generate token
	token := service.GenerateJWT(dbUser.Id, dbUser.Username, dbUser.Role)
	// sign the jwt
	tokenString, err := service.SignToken(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to sign the token"})
	}

	// set the cookie
	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		HTTPOnly: true,
		SameSite: "lax",
	}
	c.Cookie(cookie)

	/*
		// parse the token
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// no idea might check on it later
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			log.Fatal(err)
		}

		var UserRole string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// if exp, id not valid
			err := db.QueryRow("SELECT role FROM users WHERE id=?", claims["id"]).Scan(&UserRole)
			fmt.Println(UserRole)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
		c.Locals("userRole", UserRole)

		// taking the claims from token using jwt.MapClaims type??
	*/

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
}

func (r *Repository) GetUsers(c *fiber.Ctx) error {
	db := r.DB

	rows, err := db.Query("SELECT id, username, password, role FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []models.User
	var user models.User

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Role); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
