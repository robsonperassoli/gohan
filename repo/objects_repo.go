package repo

import (
	"github.com/jmoiron/sqlx"
)

type UpsertObjectParams struct {
	ID      string
	Type    string
	Display string
}

func UpsertObject(db *sqlx.DB, obj UpsertObjectParams) error {
	query := "INSERT INTO objects (id, type, display) VALUES (:id, :type, :display) ON CONFLICT DO NOTHING"
	args := map[string]interface{}{
		"id":      obj.ID,
		"type":    obj.Type,
		"display": obj.Display,
	}

	_, err := db.NamedExec(query, args)
	return err
}
