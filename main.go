package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Object struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Display string `json:"display"`
}

type Event struct {
	Timestamp time.Time `json:"timestamp"`

	Verb string `json:"verb"`

	Direct        Object `json:"direct"`
	Indirect      Object `json:"indirect"`
	Prepositional Object `json:"prepositional"`

	Context string `json:"context"`
}

func handlePostEvents(c *fiber.Ctx) error {
	events := new([]Event)
	err := c.BodyParser(&events)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	fmt.Println(events)

	c.Status(fiber.StatusCreated)
	return nil
}

func main() {
	app := fiber.New()

	app.Post("/events", handlePostEvents)

	app.Listen(":4100")
}
