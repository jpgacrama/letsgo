package mysql

import (
	"context"
	"database/sql"
	"log"
	"snippetbox/pkg/models"
	"time"
)

type SnippetDatabase struct {
	ctx             context.Context
	tx              *sql.Tx
	db              *sql.DB
	infoLog         *log.Logger
	errorLog        *log.Logger
	LatestStatement *sql.Stmt
	InsertStatement *sql.Stmt
	GetStatement    *sql.Stmt
}

// NOTE: It is now the caller's responsibility to close EACH of the Statements!
func NewSnippetModel(db *sql.DB, infolog, errorlog *log.Logger) (*SnippetDatabase, error) {
	snippetModel := &SnippetDatabase{db: db, infoLog: infolog, errorLog: errorlog}
	err := snippetModel.initializeContext()
	if err != nil {
		snippetModel.errorLog.Printf("\n\t--- InitDatabase(): Error Initializing Context: %s ---", err)
		return nil, err
	}

	// Latest() Prepared Statement
	latestStatement, err := snippetModel.tx.PrepareContext(snippetModel.ctx, `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`)
	if err != nil {
		snippetModel.errorLog.Printf("\n\t--- Latest(): Error Preparing Statement: %s ---", err)
		return nil, err
	}

	// Insert Prepared Statement
	insertStatement, err := snippetModel.tx.PrepareContext(snippetModel.ctx, `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`)
	if err != nil {
		snippetModel.errorLog.Printf("\n\t--- Get(): Error Preparing Statement: %s ---", err)
		return nil, err
	}

	// Get Prepared Statement
	getStatement, err := snippetModel.tx.PrepareContext(snippetModel.ctx, `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`)
	if err != nil {
		snippetModel.errorLog.Printf("\n\t--- Get(): Error Preparing Statement: %s ---", err)
		return nil, err
	}

	// Assigning the SQL Prepapred Statements
	snippetModel.LatestStatement = latestStatement
	snippetModel.InsertStatement = insertStatement
	snippetModel.GetStatement = getStatement
	return snippetModel, nil
}

func (m *SnippetDatabase) Close() {
	m.db.Close()
	m.LatestStatement.Close()
	m.InsertStatement.Close()
	m.GetStatement.Close()
}

// NOTE: rows.Close() must be called by the calling function!
func (m *SnippetDatabase) Latest() ([]*models.Snippet, error) {
	m.infoLog.Printf("\n\tLatest() called")
	if m.LatestStatement == nil {
		m.errorLog.Fatalf("\n\t---- Call NewSnippetModel() first----")
	}

	rows, err := m.LatestStatement.QueryContext(m.ctx)
	if err != nil {
		m.errorLog.Printf("\n\t--- Latest(): Error Querying Statement: %s ---", err)
		m.tx.Rollback()
		return nil, err
	}

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		expiresString := ""
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &expiresString)
		if err != nil {
			m.errorLog.Printf("\n\t--- Error: %s ---", err)
			return nil, err
		}
		s.Expires, err = time.Parse(time.RFC3339, expiresString)
		if err != nil {
			m.errorLog.Printf("\n\t--- Error: %s ---", err)
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

// This function takes the title, content and the time it expires
func (m *SnippetDatabase) Insert(title, content, numOfDaysToExpire string) (int, error) {
	if m.LatestStatement == nil {
		m.errorLog.Fatalf("\n\t---- Call NewSnippetModel() first----")
	}

	errorValue := -1
	// Convert expires to a string representing the number of days
	result, err := m.InsertStatement.ExecContext(m.ctx, title, content, numOfDaysToExpire)
	if err != nil {
		m.errorLog.Printf("\n\tError: %s", err)
		m.tx.Rollback()
		return errorValue, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		m.errorLog.Printf("\n\tError: %s", err)
		return errorValue, err
	}
	return int(id), nil
}

func (m *SnippetDatabase) Get(id int) (*models.Snippet, error) {
	if m.LatestStatement == nil {
		// Assumes that even the loggers for SnippetModel were not set yet
		m.errorLog.Fatalf("\n\t---- Call NewSnippetModel() first----")
	}

	s := &models.Snippet{}
	expiresString := ""
	err := m.GetStatement.QueryRowContext(m.ctx, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &expiresString)
	switch {
	case err == sql.ErrNoRows:
		m.errorLog.Printf("\n\t--- Error: %s ---", err)
		m.tx.Rollback()
		return nil, models.ErrNoRecord
	case err != nil:
		m.errorLog.Printf("\n\t--- Error: %s ---", err)
		m.tx.Rollback()
		return nil, err
	default:
		s.Expires, err = time.Parse(time.RFC3339, expiresString)
		if err != nil {
			m.errorLog.Printf("\n\t--- Error: %s ---", err)
			return nil, err
		}
		m.infoLog.Printf("ID is %v, created on %s\n", s.ID, s.Created)
		return s, nil
	}
}

func (m *SnippetDatabase) initializeContext() error {
	if m.ctx == nil {
		m.ctx = context.Background()
	}
	tx, err := m.db.BeginTx(m.ctx, nil)
	if err != nil {
		m.errorLog.Printf("\n\t--- initializeContext(): Error Beginning Transaction: %s ---", err)
		return err
	}
	m.tx = tx
	return nil
}
