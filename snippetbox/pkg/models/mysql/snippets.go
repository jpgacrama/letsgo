package mysql

import (
	"database/sql"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Close() {
	m.DB.Close()
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		println("\t--- Insert(): Error Preparing Statement ---")
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, content, expires)
	if err != nil {
		println("\t--- Insert(): Error Executing Context ---")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		println("\t--- Insert(): Getting Last Insert ID ---")
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Record, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Record{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Record, error) {
	return nil, nil
}
