package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateEventParams struct {
	Verb            string
	DirectID        string
	IndirectID      string
	PrepositionalID string
	Context         string
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  len(s) > 0,
	}
}

func InsertEvent(db *sqlx.DB, params CreateEventParams) error {
	query := `INSERT INTO events
	(id, timestamp, verb, direct_object_id, indirect_object_id, prepositional_object_id, context)
	VALUES (:id, :timestamp, :verb, :direct_object_id, :indirect_object_id, :prepositional_object_id, :context)`

	args := map[string]interface{}{
		"id":                      uuid.New().String(),
		"timestamp":               time.Now().UTC(),
		"verb":                    params.Verb,
		"direct_object_id":        params.DirectID,
		"indirect_object_id":      params.IndirectID,
		"prepositional_object_id": NewNullString(params.PrepositionalID),
		"context":                 params.Context,
	}

	_, err := db.NamedExec(query, args)

	return err
}
