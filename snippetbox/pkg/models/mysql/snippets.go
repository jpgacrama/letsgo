package mysql

import (
	"context"
	"database/sql"
	"log"
	"snippetbox/pkg/models"
	"time"
)

type SnippetDatabase struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *SnippetDatabase) Close() {
	m.DB.Close()
}

func (m *SnippetDatabase) Insert(title, content, expires string) (int, error) {
	m.InfoLog.Println("--- Inside Insert() ---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		m.ErrorLog.Println("\t--- Insert(): Error Preparing Statement ---")
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, title, content, expires)
	if err != nil {
		m.ErrorLog.Println("\t--- Insert(): Error Executing Context ---")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		m.ErrorLog.Println("\t--- Insert(): Error Getting Last Insert ID ---")
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetDatabase) Get(id int) (*models.Snippet, error) {
	m.InfoLog.Println("--- Inside Get() ---")
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		m.ErrorLog.Println("\t--- Get(): No Record ---")
		return nil, models.ErrNoRecord
	} else if err != nil {
		m.ErrorLog.Println("\t--- Get(): Error Scanning ---")
		return nil, err
	}
	return s, nil
}

func (m *SnippetDatabase) FindByID(id string) (*models.Snippet, error) {
	user := new(models.Snippet)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(
		ctx, "SELECT id, title, content, created, expires FROM snippets WHERE id = ?", id).Scan(
		&user.ID, &user.Title, &user.Content, &user.Created, &user.Expires)
	if err != nil {
		m.ErrorLog.Println("\t--- FindByID(): Error Querying by ID ---")
		return nil, err
	}
	return user, nil
}

func (m *SnippetDatabase) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
