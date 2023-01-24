package snippetbox_test

import (
	"database/sql"
	"log"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"strconv"
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
	db, mock := NewMock()
	infoLog, errorLog := server.CreateLoggers()
	repo := &mysql.SnippetDatabase{
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog}
	defer func() {
		repo.Close()
	}()
	t.Run("Insert OK Case", func(t *testing.T) {
		query := "INSERT INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(
			u.Title,
			u.Content,
			u.Expires).WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := repo.Insert(u.Title, u.Content, u.Expires)
		assert.NoError(t, err)
	})
	t.Run("Insert NOK Case", func(t *testing.T) {
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

func TestGet(t *testing.T) {
	db, mock := NewMock()
	infoLog, errorLog := server.CreateLoggers()
	repo := &mysql.SnippetDatabase{
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog}
	defer func() {
		repo.Close()
	}()
	t.Run("Get OK Case", func(t *testing.T) {
		query := "SELECT id, title, content, created, expires FROM snippets WHERE expires \\> UTC_TIMESTAMP\\(\\) AND id = \\?"
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

		user, err := repo.FindByID(strconv.Itoa(u.ID))
		assert.Empty(t, user)
		assert.Error(t, err)
	})
}
