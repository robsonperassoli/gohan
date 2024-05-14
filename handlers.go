package main

import (
	"fmt"
	"gohan/repo"
	"gohan/views"
	"log"
	"log/slog"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jmoiron/sqlx"
)

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CreateEventObject struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Display string `json:"display"`
}

type CreateEvent struct {
	Subject CreateEventObject `json:"subject"`
	Verb    string            `json:"verb"`

	Direct        CreateEventObject `json:"direct"`
	Indirect      CreateEventObject `json:"indirect"`
	Prepositional CreateEventObject `json:"prepositional"`
	Context       string            `json:"context"`
}

func HandleError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
		Success: false,
		Message: err.Error(),
	})
}

func HandleHome(c *fiber.Ctx, db *sqlx.DB) error {
	objectIdsParam := c.Query("object_ids", "")
	verb := c.Query("verb", "")

	objectIds := []string{}
	if objectIdsParam != "" {
		objectIds = strings.Split(objectIdsParam, ",")

	}

	selectedObjects := []repo.Object{}
	if len(objectIds) > 0 {
		var err error
		selectedObjects, err = repo.GetObjectByIds(db, objectIds)
		if err != nil {
			return err
		}
		fmt.Println(selectedObjects)
	}

	events, err := repo.ListEvents(db, repo.ListFilters{
		ObjectIDs: objectIds,
		Verb:      verb,
	})
	if err != nil {
		slog.Error("Could not list events", err)
	}
	fmt.Println(selectedObjects)

	handler := adaptor.HTTPHandler(templ.Handler(views.Home(events, verb, selectedObjects)))
	return handler(c)
}

func HandlePostEvents(c *fiber.Ctx, db *sqlx.DB) error {
	events := new([]CreateEvent)
	err := c.BodyParser(&events)
	if err != nil {
		return err
	}

	for _, e := range *events {
		err = repo.UpsertObject(db, repo.UpsertObjectParams{
			ID:      e.Subject.ID,
			Type:    e.Subject.Type,
			Display: e.Subject.Display,
		})
		if err != nil {
			log.Fatal("Could not save subject", err)
		}

		err = repo.UpsertObject(db, repo.UpsertObjectParams{
			ID:      e.Direct.ID,
			Type:    e.Direct.Type,
			Display: e.Direct.Display,
		})
		if err != nil {
			log.Fatal("Could not save direct object", err)
		}

		if e.Indirect.ID != "" {
			err = repo.UpsertObject(db, repo.UpsertObjectParams{
				ID:      e.Indirect.ID,
				Type:    e.Indirect.Type,
				Display: e.Indirect.Display,
			})
			if err != nil {
				log.Fatal("Could not save indirect object", err)
			}
		}

		if e.Prepositional.ID != "" {
			err = repo.UpsertObject(db, repo.UpsertObjectParams{
				ID:      e.Prepositional.ID,
				Type:    e.Prepositional.Type,
				Display: e.Prepositional.Display,
			})
			if err != nil {
				log.Fatal("Could not save prepositional object", err)
			}
		}

		err = repo.InsertEvent(db, repo.CreateEventParams{
			Verb:            e.Verb,
			SubjectID:       e.Subject.ID,
			DirectID:        e.Direct.ID,
			IndirectID:      e.Indirect.ID,
			PrepositionalID: e.Prepositional.ID,
			Context:         e.Context,
		})
		if err != nil {
			log.Fatal("Could no Save the event", err)
		}
	}

	c.Status(fiber.StatusCreated)
	return nil
}

func HandleGetObjectById(c *fiber.Ctx, db *sqlx.DB) error {
	id := c.Params("id")

	obj, err := repo.GetObjectById(db, id)
	if err != nil {
		slog.Error("Could not fetch object", err)
	}

	handler := adaptor.HTTPHandler(templ.Handler(views.Object(obj)))
	return handler(c)
}
