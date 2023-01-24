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

var u = &models.Record{
	ID:      1,
	Title:   "Title",
	Content: "Content",
	Created: time.Now(),
	Expires: "1",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestInsert(t *testing.T) {
	t.Run("Insert OK Case", func(t *testing.T) {
		db, mock := NewMock()
		repo := &mysql.SnippetModel{DB: db}
		defer func() {
			repo.Close()
		}()

		query := "INSERT OR UPDATE INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(
			u.Title,
			u.Content,
			u.Expires).WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := repo.Insert(u.Title, u.Content, u.Expires)
		assert.NoError(t, err)
	})
	t.Run("Insert NOK Case", func(t *testing.T) {
		db, mock := NewMock()
		repo := &mysql.SnippetModel{DB: db}
		defer func() {
			repo.Close()
		}()

		query := "INSERT OR UPDATE INTO snippet \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(
			u.Title,
			u.Content,
			u.Expires).WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := repo.Insert(u.Title, u.Content, u.Expires)
		assert.Error(t, err)
	})
}
