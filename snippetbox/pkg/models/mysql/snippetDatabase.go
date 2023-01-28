package mysql

import (
	"context"
	"database/sql"
	"log"
	"snippetbox/pkg/models"
)

type SnippetDatabase struct {
	ctx      context.Context
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *SnippetDatabase) Close() {
	m.DB.Close()
}

func (m *SnippetDatabase) Latest() ([]*models.Snippet, error) {
	m.InfoLog.Println("--- Inside Latest() ---")
	m.initializeContext()
	stmt, err := m.DB.Prepare(`SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`)
	if err != nil {
		m.ErrorLog.Printf("\n\t--- Latest(): Error Preparing Statement: %s ---", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(m.ctx)
	if err != nil {
		m.ErrorLog.Printf("\n\t--- Latest(): Error Querying Statement: %s ---", err)
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			m.ErrorLog.Printf("\n\t--- Latest(): Error Scanning: %s ---", err)
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetDatabase) Insert(title, content, expires string) (int, error) {
	m.InfoLog.Println("--- Inside Insert() ---")
	m.initializeContext()

	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		m.ErrorLog.Printf("\n\t--- Insert(): Error Preparing Statement: %s ---", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(m.ctx, title, content, expires)
	if err != nil {
		m.ErrorLog.Println("\n\t--- Insert(): Error Executing Context ---")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		m.ErrorLog.Println("\n\t--- Insert(): Error Getting Last Insert ID ---")
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetDatabase) Get(id int) (*models.Snippet, error) {
	m.InfoLog.Println("--- Inside Get() ---")
	m.initializeContext()
	stmt, err := m.DB.Prepare(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`)
	if err != nil {
		m.ErrorLog.Printf("\n\t--- Get(): Error Preparing Statement: %s ---", err)
		return nil, err
	}
	defer stmt.Close()

	s := &models.Snippet{}
	err = stmt.QueryRowContext(m.ctx, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	switch {
	case err == sql.ErrNoRows:
		m.ErrorLog.Println("\n\t--- Get(): No Record ---")
		return nil, models.ErrNoRecord
	case err != nil:
		m.ErrorLog.Print("\n\t--- Get(): Error Querying:", err, " ---")
		return nil, err
	default:
		log.Printf("ID is %v, created on %s\n", s.ID, s.Created)
		return s, nil
	}
}

func (m *SnippetDatabase) initializeContext() {
	if m.ctx == nil {
		m.ctx = context.Background()
	}
}
