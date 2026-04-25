package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Get("/", HelloWorld)

	panic(app.Listen(":" + os.Getenv("PORT")))
}

func HelloWorld(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": "true",
	})
}
