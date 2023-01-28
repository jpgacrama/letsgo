package snippetbox_test

import (
	"database/sql"
	"log"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var sampleDatabaseContent = &models.Snippet{
	ID:      0,
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
			sampleDatabaseContent.Title,
			sampleDatabaseContent.Content,
			sampleDatabaseContent.Expires).WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := repo.Insert(sampleDatabaseContent.Title, sampleDatabaseContent.Content, sampleDatabaseContent.Expires)
		assert.NoError(t, err)
	})
	t.Run("Insert NOK Case", func(t *testing.T) {
		query := "INSERT OR UPDATE INTO snippet \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(
			sampleDatabaseContent.Title,
			sampleDatabaseContent.Content,
			sampleDatabaseContent.Expires).WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := repo.Insert(sampleDatabaseContent.Title, sampleDatabaseContent.Content, sampleDatabaseContent.Expires)
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
		query := "SELECT ..."
		prep := mock.ExpectPrepare(query)
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "1")
		prep.ExpectQuery().WithArgs(sampleDatabaseContent.ID).WillReturnRows(rows)

		output, err := repo.Get(sampleDatabaseContent.ID)
		assert.NotNil(t, output)
		assert.NoError(t, err)
	})
	t.Run("Get NOK Case", func(t *testing.T) {
		query := "SELECT ..."
		prep := mock.ExpectPrepare(query)
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "1")

		wrongId := 2
		output, err := repo.Get(wrongId)
		assert.Nil(t, output)
		prep.ExpectQuery().WithArgs().WillReturnError(err)
		assert.Error(t, err)
	})
}
