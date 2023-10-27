package main

import "github.com/gofiber/fiber/v2"

func main(){
	var app *fiber.App = fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from FIBER!")
	})

	app.Listen(":3000")
}