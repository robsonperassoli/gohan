package repo

import (
	"fmt"

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
	objs, err := GetObjectByIds(db, []string{id})
	if err != nil {
		return Object{}, err
	}

	return objs[0], nil
}

func GetObjectByIds(db *sqlx.DB, ids []string) ([]Object, error) {
	objs := []Object{}

	idArgs := []interface{}{}
	for _, str := range ids {
		idArgs = append(idArgs, str)
	}

	query := "SELECT id, type, display FROM objects where id IN (?)"
	query, args, err := sqlx.In(query, idArgs...)
	if err != nil {
		return objs, err
	}

	fmt.Println(args)
	query = db.Rebind(query)
	if err := db.Select(&objs, query, args...); err != nil {
		return objs, err
	}

	return objs, nil
}
