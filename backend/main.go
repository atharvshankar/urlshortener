package main

import (
	"fmt"
	"strings"
	"github.com/gofiber/fiber/v2"
)


func convertToBase62(counter int64) string {
	if counter == 0 {
		return "0"
	}
	base62Helper := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for counter>0 {
		digit := counter%62
		sb.WriteString(string(base62Helper[digit]))
		counter = counter/62
	}
	return sb.String()
}

func main(){

	var counter int64
	counter = 100000

	//Read next counter value from DB during startup. Helpful when backend crashes or is restarted.
	//TODO: counter = ReadCounterValueFromDB()
	var app *fiber.App = fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(convertToBase62(counter))
		counter++
		//Testing the ConvertToBase62Function
		return c.SendString("Hello from FIBER!")
	})

	//DB : urlshortenerdb TABLE: urlmapping
	//TODO: Implement the shorten API which takes in the long URL in JSON Body of the request and returns the short URL
	app.Post("/shorten/", func(c *fiber.Ctx) error {
		// counter++
		// return c.SendString(convertToBase62(counter))
	})

	//FIXME: Persist the shortedned URLs into DB, and fetch the corresponding long URLs from DB when user requests for it.
	app.Get("/:shortenedurl/", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("shortenedurl"))
		if c.Params("shortenedurl") == "V0q"{
			return c.Redirect("https://www.google.co.in")
		}
		return c.Redirect("http://localhost:3000/")
	})

	app.Listen(":3000")
}