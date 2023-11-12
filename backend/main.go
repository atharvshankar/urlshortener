package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"time"
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
	const domainName string = "http://localhost:3000/"
	
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
	// counter = 100000
	//Read next counter value from DB during startup. Helpful when backend crashes or is restarted.
	// counter = ReadCounterValueFromDB()
	row := db.QueryRow("SELECT value FROM tracker WHERE name = $1","counter")
	switch err := row.Scan(&counter); err{
	case sql.ErrNoRows:
		//Stop backend if you cant fetch the counter value.
		os.Exit(1)
	}


	
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

	type response struct {
		LongURL string `json:"long_url"`
		ShortURL string `json:"short_url"`
	}

	//DB : urlshortenerdb TABLE: urlmapping
	app.Post("/shorten/", func(c *fiber.Ctx) error {
		url := new(URL)
		if err := c.BodyParser(url); err != nil {
			fmt.Println(err);
            return c.SendString("Unable to process request!")
        }
		fmt.Println(url.LongURL);

		//Ignoring domain, only inserting the unique part of shorturl inside db.
		shortURLHash := convertToBase62(counter)
		counter += 1
		db.Exec("UPDATE tracker SET value=$1 WHERE name=$2",counter,"counter")
		insertSQLStatement:= `INSERT INTO urlmapping(long_url, short_url, created_on) VALUES ($1, $2, $3)`
		_, err := db.Exec(insertSQLStatement, url.LongURL, shortURLHash, time.Now())
		if err!= nil {
			panic(err)
		}
		completeShortUrl := domainName + shortURLHash
		reponseObj := response{
			LongURL: url.LongURL,
			ShortURL: completeShortUrl,
		}
		return c.JSON(reponseObj)
	})

	app.Get("/:shortenedurl/", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("shortenedurl"))
		sqlStatement := `SELECT long_url from urlmapping WHERE short_url = $1;`
		var longURL string
		row := db.QueryRow(sqlStatement,c.Params("shortenedurl"))
		switch err := row.Scan(&longURL); err{
		case sql.ErrNoRows:
			//FIXME: Add 404 not found page at the error URL
			return c.Redirect("http://localhost:3000/error")
		case nil:
			return c.Redirect(longURL)
		default:
			panic(err)
		}
	})

	app.Listen(":3000")
}
