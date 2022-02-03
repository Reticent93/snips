package postgres

import (
	"database/sql"
	"errors"
	"github.com/Reticent93/snips/pkg/models"
	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

type SnipModel struct {
	DB *bun.DB
}

func (m *SnipModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippet(title, content, created, expires)
VALUES(?, ?, CURRENT_TIMESTAMP, CURRENT_DATE + INTERVAL ? DAY) RETURNING id`

	id := 0
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnipModel) Get(id int) (*models.Snip, error) {
	s := &models.Snip{}
	stmnt := `SELECT id, title, content, created, expires FROM snippet
			WHERE expires > current_date AND id = ?`

	err := m.DB.QueryRow(stmnt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *SnipModel) Latest() ([]*models.Snip, error) {
	stmt := `SELECT id, title, content,created,expires FROM snippet
	WHERE expires > CURRENT_TIMESTAMP ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	snippets := []*models.Snip{}

	for rows.Next() {
		s := &models.Snip{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return snippets, nil

}
