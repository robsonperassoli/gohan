package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/vrischmann/envconfig"

	_ "github.com/lib/pq"
)

var conf struct {
	DB struct {
		URL string `envconfig:"default=postgresql://postgres:postgres@localhost/gohan?sslmode=disable"`
	}
	PORT string `envconfig:"default=4100"`
}

func MustConnectDB(connString string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

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
	app.Static("/", "./public")
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return HandleHome(c, db)
	})

	app.Post("/events", func(c *fiber.Ctx) error {
		return HandlePostEvents(c, db)
	})

	app.Get("/objects/:id", func(c *fiber.Ctx) error {
		return HandleGetObjectById(c, db)
	})

	addr := fmt.Sprintf(":%s", conf.PORT)
	log.Fatal(app.Listen(addr))
}
