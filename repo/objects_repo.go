package repo

import (
	"github.com/jmoiron/sqlx"
)

type UpsertObjectParams struct {
	ID      string
	Type    string
	Display string
}

type Object struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Display string `json:"display"`
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

func GetObjectById(db *sqlx.DB, id string) (Object, error) {
	query := "SELECT id, type, display FROM objects where id = $1"
	obj := Object{}
	if err := db.Get(&obj, query, id); err != nil {
		return obj, err
	}

	return obj, nil
}
