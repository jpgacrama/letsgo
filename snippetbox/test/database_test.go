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

var created = time.Now()
var sampleDatabaseContent = &models.Snippet{
	ID:      1,
	Title:   "Title",
	Content: "Content",
	Created: created,
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

	// New mocks due to NewSnippetModel() factory
	mock.ExpectBegin()
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for Latest Statement
	query := "INSERT INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"
	prep := mock.ExpectPrepare(query)
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items

	repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	defer func() {
		if err == nil {
			repo.Close()
		}
	}()

	if err != nil {
		log.Fatalf("Creating NewSnippetModel failed")
		return
	}
	t.Run("Insert OK Case", func(t *testing.T) {
		prep.ExpectExec().WithArgs(
			sampleDatabaseContent.Title,
			sampleDatabaseContent.Content,
			"1").WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := repo.Insert(sampleDatabaseContent.Title, sampleDatabaseContent.Content, "1")
		assert.NoError(t, err)
	})
	t.Run("Insert NOK Case", func(t *testing.T) {
		query := "INSERT INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"
		mock.ExpectQuery(query).WithArgs(
			sampleDatabaseContent.Title,
			sampleDatabaseContent.Content,
			"1").WillReturnError(err)
		_, err := repo.Insert(sampleDatabaseContent.Title, sampleDatabaseContent.Content, "1")
		assert.Error(t, err)
	})
}

func TestGet(t *testing.T) {
	db, mock := NewMock()
	infoLog, errorLog := server.CreateLoggers()

	// New mocks due to NewSnippetModel() factory
	mock.ExpectBegin()
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for Latest Statement
	_ = mock.ExpectPrepare("INSERT ...")

	query := "SELECT id, title, content, created, expires FROM snippets WHERE expires \\> UTC_TIMESTAMP\\(\\) AND id \\= \\?"
	prep := mock.ExpectPrepare(query) // SELECT for just one of the items

	repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	defer func() {
		if err == nil {
			repo.Close()
		}
	}()

	if err != nil {
		log.Fatalf("Creating NewSnippetModel failed")
		return
	}

	t.Run("Get() OK Case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(sampleDatabaseContent.ID).WillReturnRows(rows)

		output, err := repo.Get(sampleDatabaseContent.ID)
		assert.NotNil(t, output)
		assert.NoError(t, err)
	})
	t.Run("Get() NOK Case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")

		wrongId := 2
		output, err := repo.Get(wrongId)
		assert.Nil(t, output)
		prep.ExpectQuery().WithArgs().WillReturnError(err)
		assert.Error(t, err)
	})
}

func TestLatest(t *testing.T) {
	db, mock := NewMock()
	infoLog, errorLog := server.CreateLoggers()

	// New mocks due to NewSnippetModel() factory
	mock.ExpectBegin()

	// SELECT for Latest Statement
	query := "SELECT id, title, content, created, expires FROM snippets WHERE expires \\> UTC_TIMESTAMP\\(\\) ORDER BY created DESC LIMIT 10"
	prep := mock.ExpectPrepare(query)
	_ = mock.ExpectPrepare("INSERT ...")
	_ = mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items

	repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
	defer func() {
		if err == nil {
			repo.Close()
		}
	}()

	if err != nil {
		log.Fatalf("Creating NewSnippetModel failed")
		return
	}
	t.Run("Latest() OK Case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WillReturnRows(rows)

		output, err := repo.Latest()
		assert.NotNil(t, output)
		assert.NoError(t, err)
	})
	t.Run("Latest() NOK Case", func(t *testing.T) {
		output, err := repo.Latest()
		prep.ExpectQuery().WillReturnError(err)
		assert.Nil(t, output)
		assert.Error(t, err)
	})
}