package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/vrischmann/envconfig"

	_ "github.com/lib/pq"
)

var conf struct {
	DB struct {
		URL string `envconfig:"default=postgresql://postgres:postgres@localhost/gohan?sslmode=disable"`
	}
}

func MustConnectDB(connString string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	return db
}

func main() {
	if err := envconfig.Init(&conf); err != nil {
		log.Fatalln(err)
	}

	db := MustConnectDB(conf.DB.URL)

	app := fiber.New(fiber.Config{
		ErrorHandler: HandleError,
	})

	app.Post("/events", func(c *fiber.Ctx) error {
		return HandlePostEvents(c, db)
	})

	log.Fatal(app.Listen(":4100"))
}
