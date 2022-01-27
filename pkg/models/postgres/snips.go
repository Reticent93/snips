package postgres

import (
	"database/sql"
	"errors"
	"github.com/Reticent93/snips/pkg/models"
	"github.com/uptrace/bun"
)

type SnipModel struct {
	DB *bun.DB
}

func (m *SnipModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippet(title, content, created, expires)
VALUES(?, ?, CURRENT_TIMESTAMP, TIMESTAMP(?, ?)) RETURNING id`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *SnipModel) Get(id int) (*models.Snip, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippet
	WHERE expires > CURRENT_TIMESTAMP AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snip{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

	//	var s = &Snip{}
	//	err := m.DB.QueryRow(`SELECT title, content, created, expires FROM snippets
	//WHERE expires > CURRENT_TIMESTAMP`, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	//
	//
	//	if err != nil {
	//		if errors.Is(err, sql.ErrNoRows) {
	//			return nil, ErrNoRecord
	//		} else {
	//			return nil, err
	//		}
	//	}
	//	return s, nil
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

	var snippets []*models.Snip

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
