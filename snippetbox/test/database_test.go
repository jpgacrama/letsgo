package snippetbox_test

import (
	"database/sql"
	"log"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var u = &models.Snippet{
	ID:      1,
	Title:   "Title",
	Content: "Content",
	Created: time.Now(),
	Expires: time.Now().AddDate(0, 0, 1),
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestInsert(t *testing.T) {
	db, mock := NewMock()
	repo := &mysql.SnippetModel{DB: db}
	defer func() {
		repo.Close()
	}()

	query := "INSERT INTO users \\(id, title, content, created, expires\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		u.ID,
		u.Title,
		u.Content,
		u.Created,
		u.Expires).WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := repo.Insert(u.Title, u.Content, u.Expires.String())
	assert.NoError(t, err)
}
