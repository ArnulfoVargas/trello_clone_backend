package routes

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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

	g.Post("/login", r.login)
	g.Post("/register", r.register)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *Router) login(c fiber.Ctx) error {
	var login models_auth.Login

	err := json.Unmarshal(c.Body(), &login)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := context.Background()
	u, err := gorm.G[models_auth.User](r.GetConn()).Where("email LIKE ?", login.Email).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(createError(defaultLoginError))
		}
		return c.SendStatus(fiber.StatusBadGateway)
	}

	if !CheckPassword(login.Password, u.Password) {
		return c.JSON(createError(defaultLoginError))
	}

	token, err := GenerateToken(u.ID)
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

func (r *Router) register(c fiber.Ctx) error {
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

	hash, err := HashPassword(register.Password)

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

	token, err := GenerateToken(user.ID)
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


func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 días
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}