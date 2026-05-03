package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"
	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/utils"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultLoginError    = "Incorrect Email or Password"
	defaultRegisterError = "Invalid Email"
	emailAlreadyUsed     = "Email already used"
)

func createError(text string) fiber.Map {
	return fiber.Map{
		"error": text,
	}
}

func (r *Router) BindAuthRoutes() {
	g := r.app.Group("/auth")

	g.Post("/login", r.genLogin("email LIKE ?"))
	g.Post("/register", r.register)
	g.Post("/slogin", r.genLogin("email LIKE ? AND role = 1"))
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *Router) genLogin(findQuery string) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		var login models_auth.Login

		err := json.Unmarshal(c.Body(), &login)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		ctx := context.Background()

		var u models_auth.User

		u, err = gorm.G[models_auth.User](r.GetConn()).Where(findQuery, login.Email).First(ctx)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(createError(defaultLoginError))
			}
			return c.SendStatus(fiber.StatusBadGateway)
		}

		if !checkPassword(login.Password, u.Password) {
			return c.JSON(createError(defaultLoginError))
		}

		token, err := utils.GenerateToken(u.ID)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
			Expires:  time.Now().Add(time.Hour * 24 * 7),
		})

		return c.JSON(fiber.Map{"success": true, "token": token})
	}
}

func (r *Router) register(c fiber.Ctx) error {

	fmt.Println("Hit register")

	var register models_auth.Register
	err := json.Unmarshal(c.Body(), &register)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !register.ValidateEmail() {
		return c.JSON(createError(defaultRegisterError))
	}

	if r, err := register.ValidatePassword(); !r {
		return c.JSON(createError(err))
	}

	if r, err := register.ValidateUsername(); !r {
		return c.JSON(createError(err))
	}

	ctx := context.Background()
	_, err = gorm.G[models_auth.User](r.GetConn()).Where("email LIKE ?", register.Email).First(ctx)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(createError(emailAlreadyUsed))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	hash, err := hashPassword(register.Password)

	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	createCtx := context.Background()
	user := &models_auth.User{
		Name:     register.Username,
		Email:    register.Email,
		Password: hash,
	}

	err = gorm.G[models_auth.User](r.GetConn()).Create(createCtx, user)

	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	})

	return c.JSON(fiber.Map{"success": true, "token": token})
}


