package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateEventParams struct {
	SubjectID       string
	Verb            string
	DirectID        string
	IndirectID      string
	PrepositionalID string
	Context         string
}

type Event struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	Subject       Object    `json:"subject"`
	Verb          string    `json:"verb"`
	Direct        Object    `json:"direct"`
	Indirect      *Object   `json:"indirect"`
	Prepositional *Object   `json:"prepositional"`
	Context       string    `json:"context"`
}

type ListFilters struct {
	ObjectIDs []string
	Verb      string
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  len(s) > 0,
	}
}

func InsertEvent(db *sqlx.DB, params CreateEventParams) error {
	query := `INSERT INTO events
	(id, timestamp, verb, subject_id, direct_object_id, indirect_object_id, prepositional_object_id, context)
	VALUES (:id, :timestamp, :verb, :subject_id, :direct_object_id, :indirect_object_id, :prepositional_object_id, :context)`

	args := map[string]interface{}{
		"id":                      uuid.New().String(),
		"timestamp":               time.Now().UTC(),
		"verb":                    params.Verb,
		"subject_id":              params.SubjectID,
		"direct_object_id":        params.DirectID,
		"indirect_object_id":      NewNullString(params.IndirectID),
		"prepositional_object_id": NewNullString(params.PrepositionalID),
		"context":                 params.Context,
	}

	_, err := db.NamedExec(query, args)

	return err
}

func ListEvents(db *sqlx.DB, filters ListFilters) ([]Event, error) {
	where := "where true"
	args := []interface{}{}

	if filters.Verb != "" {
		where += " and verb = ?"
		args = append(args, filters.Verb)
	}

	if len(filters.ObjectIDs) > 0 {
		where += " and (direct_object_id in (?) or indirect_object_id in (?) or prepositional_object_id in (?) or subject_id in (?))"
		args = append(args, filters.ObjectIDs, filters.ObjectIDs, filters.ObjectIDs, filters.ObjectIDs)
	}

	query := fmt.Sprintf(`SELECT
		e.id, e.timestamp, e.verb, e.context,
		subject.id, subject.type, subject.display,
		direct.id, direct.type, direct.display,
		indirect.id, indirect.type, indirect.display,
		prepositional.id, prepositional.type, prepositional.display
	FROM
		events e
	INNER JOIN objects subject on subject.id = e.subject_id
	INNER JOIN objects direct on direct.id = e.direct_object_id
	LEFT JOIN objects indirect on indirect.id = e.indirect_object_id
	LEFT JOIN objects prepositional on prepositional.id = e.prepositional_object_id
	%s
	ORDER BY
		e.timestamp desc`, where)

	query, args, err := sqlx.In(query, args...)
	query = db.Rebind(query)
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	var events []Event

	for rows.Next() {
		var (
			id                    string
			timestamp             time.Time
			verb                  string
			context               string
			subject_id            string
			subject_type          string
			subject_display       string
			direct_id             string
			direct_type           string
			direct_display        string
			indirect_id           sql.NullString
			indirect_type         sql.NullString
			indirect_display      sql.NullString
			prepositional_id      sql.NullString
			prepositional_type    sql.NullString
			prepositional_display sql.NullString
		)

		err = rows.Scan(&id, &timestamp, &verb, &context, &subject_id, &subject_type, &subject_display,
			&direct_id, &direct_type, &direct_display, &indirect_id, &indirect_type, &indirect_display,
			&prepositional_id, &prepositional_type, &prepositional_display)

		subject := Object{
			ID:      subject_id,
			Type:    subject_type,
			Display: subject_display,
		}

		direct := Object{
			ID:      direct_id,
			Type:    direct_type,
			Display: direct_display,
		}

		var indirect *Object
		if indirect_id.Valid {
			indirect = &Object{
				ID:      indirect_id.String,
				Type:    indirect_type.String,
				Display: indirect_display.String,
			}
		}

		var prepositional *Object
		if prepositional_id.Valid {
			prepositional = &Object{
				ID:      prepositional_id.String,
				Type:    prepositional_type.String,
				Display: prepositional_display.String,
			}
		}

		event := Event{
			ID:            id,
			Timestamp:     timestamp,
			Verb:          verb,
			Context:       context,
			Subject:       subject,
			Direct:        direct,
			Indirect:      indirect,
			Prepositional: prepositional,
		}

		events = append(events, event)
	}

	return events, nil
}
