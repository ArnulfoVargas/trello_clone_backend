package routes

import (
	"context"

	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/middlewares"
	models_auth "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/auth"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func (r *Router) BindUserRoutes() {
	g := r.app.Group("/user", middlewares.AuthMiddleware)

	g.Get("/", r.getUserData)
}

func (r *Router) getUserData(c fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    
    ctx := context.Background()
    u, err := gorm.G[models_auth.User](r.GetConn()).Where("id = ?", userID).First(ctx)
    if err != nil {
        return c.SendStatus(fiber.StatusNotFound)
    }

    return c.JSON(fiber.Map{
        "name":  u.Name,
        "email": u.Email,
    })
}
