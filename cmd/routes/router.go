package routes

import (
	"os"

	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

type Router struct {
	app *fiber.App
	db *database.Database
}

func NewRouter() *Router {
	return &Router{
		app: fiber.New(),
	}
}

func (r *Router) SetDb(db *database.Database) {
	r.db = db
}

func (r *Router) ServeHttp() {
	panic(r.app.Listen(":" + os.Getenv("PORT")))
}

func (r *Router) GetConn() (*gorm.DB) {
	return r.db.Db
}

func (r *Router) SetCors() {
	r.app.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowCredentials: true,
}))
}