package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	_"github.com/lib/pq"
)





func convertToBase62(counter int64) string {
	if counter == 0 {
		return "0"
	}
	base62Helper := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for counter > 0 {
		digit := counter % 62
		sb.WriteString(string(base62Helper[digit]))
		counter = counter / 62
	}
	return sb.String()
}

func main() {
	godotenv.Load(".env")
	var (
		host string = os.Getenv("DB_HOST")
		port, err = strconv.Atoi(os.Getenv("DB_PORT"))
		user string = os.Getenv("DB_USER")
		password string = os.Getenv("DB_PASSWORD")
		dbname string = os.Getenv("DB_NAME")
	)
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres",psqlInfo)

	if err != nil{
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil{
		panic(err)
	}

	fmt.Println("Connected to DB Successfully!")

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


	type URL struct {
		LongURL string `json:"url" form:"url"`
	}

	//DB : urlshortenerdb TABLE: urlmapping
	//TODO: Implement the shorten API which takes in the long URL in JSON Body of the request and returns the short URL
	app.Post("/shorten/", func(c *fiber.Ctx) error {
		// counter++
		// return c.SendString(convertToBase62(counter))
		url := new(URL)
		if err := c.BodyParser(url); err != nil {
			fmt.Println(err);
            return c.SendString("Unable to process request!")
        }
		fmt.Println(url.LongURL);
		return c.SendString("You want to shorten: " + url.LongURL)
	})

	//FIXME: Persist the shortedned URLs into DB, and fetch the corresponding long URLs from DB when user requests for it.
	// app.Get("/:shortenedurl/", func(c *fiber.Ctx) error {
	// 	fmt.Println(c.Params("shortenedurl"))
	// 	if c.Params("shortenedurl") == "V0q" {
	// 		return c.Redirect("https://www.google.co.in")
	// 	}
	// 	return c.Redirect("http://localhost:3000/")
	// })

	app.Listen(":3000")
}
