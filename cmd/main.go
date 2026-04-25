package main

import (
	"os"

	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/database"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	var conn = &database.Database{}
	conn.ConnectDB(os.Getenv("DB_STRING"))

	app := fiber.New()

	app.Get("/", HelloWorld)

	panic(app.Listen(":" + os.Getenv("PORT")))
}

func HelloWorld(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": "true",
	})
}
