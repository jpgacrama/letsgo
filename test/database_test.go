package test

import (
	"database/sql"
	"log"
	"snippetbox/cmd/server"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
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
		log.Printf("Creating NewSnippetModel failed")
		return
	}
	t.Run("Insert OK Case", func(t *testing.T) {
		prep.ExpectExec().WithArgs(
			"Title",
			"Content",
			"1").WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := repo.Insert("Title", "Content", "1")
		assert.NoError(t, err)
	})
	t.Run("Insert NOK Case", func(t *testing.T) {
		query := "INSERT INTO snippets \\(title, content, created, expires\\) VALUES\\(\\?, \\?, UTC_TIMESTAMP\\(\\), DATE_ADD\\(UTC_TIMESTAMP\\(\\), INTERVAL \\? DAY\\)\\)"
		mock.ExpectQuery(query).WithArgs(
			"Title",
			"Content",
			"1").WillReturnError(err)
		_, err := repo.Insert("Title", "Content", "1")
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
		log.Printf("Creating NewSnippetModel failed")
		return
	}

	t.Run("Get() OK Case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WithArgs(0).WillReturnRows(rows)

		output, err := repo.Get(0)
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
	t.Run("Latest() OK Case", func(t *testing.T) {
		db, mock := NewMock()
		infoLog, errorLog := server.CreateLoggers()

		// New mocks due to NewSnippetModel() factory
		mock.ExpectBegin()

		// SELECT for Latest Statement
		query := "SELECT id, title, content, created, expires FROM snippets WHERE expires \\> UTC_TIMESTAMP\\(\\) ORDER BY created DESC LIMIT 10"
		prep := mock.ExpectPrepare(query)
		_ = mock.ExpectPrepare("INSERT ...")
		_ = mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items		repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)

		repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
		defer func() {
			if err == nil {
				repo.Close()
			}
		}()

		defer func() {
			if err == nil {
				repo.Close()
			}
		}()

		if err != nil {
			log.Printf("Creating NewSnippetModel failed")
			return
		}

		rows := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"})
		rows.AddRow(0, "Title", "Content", time.Now(), "2024-01-24T10:23:42Z")
		prep.ExpectQuery().WillReturnRows(rows)

		output, err := repo.Latest()
		assert.NotNil(t, output)
		assert.NoError(t, err)
	})
	t.Run("Latest() NOK Case - No Records found", func(t *testing.T) {
		db, mock := NewMock()
		infoLog, errorLog := server.CreateLoggers()

		// New mocks due to NewSnippetModel() factory
		mock.ExpectBegin()

		// SELECT for Latest Statement
		query := "SELECT id, title, content, created, expires FROM snippets WHERE expires \\> UTC_TIMESTAMP\\(\\) ORDER BY created DESC LIMIT 10"
		prep := mock.ExpectPrepare(query)
		_ = mock.ExpectPrepare("INSERT ...")
		_ = mock.ExpectPrepare("SELECT ...") // SELECT for just one of the items		repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)

		repo, err := mysql.NewSnippetModel(db, infoLog, errorLog)
		defer func() {
			if err == nil {
				repo.Close()
			}
		}()

		if err != nil {
			log.Printf("Creating NewSnippetModel failed")
			return
		}
		prep.ExpectQuery().WillReturnError(err)
		output, err := repo.Latest()
		assert.Nil(t, output)
		assert.Error(t, err)
	})
	t.Run("Latest() NOK Case - No Prepared Query exists", func(t *testing.T) {
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
			log.Printf("Creating NewSnippetModel failed")
			return
		}
		repo.LatestStatement = nil
		output, err := repo.Latest()
		prep.ExpectQuery().WillReturnError(err)
		assert.Nil(t, output)
		assert.Error(t, err)
	})
	t.Run("Insert() NOK Case - No Prepared Query exists", func(t *testing.T) {
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
			log.Printf("Creating NewSnippetModel failed")
			return
		}
		repo.InsertStatement = nil
		output, err := repo.Insert("Title", "Content", "1")
		prep.ExpectQuery().WillReturnError(err)
		assert.EqualValues(t, -1, output)
		assert.Error(t, err)
	})
	t.Run("Get() NOK Case - No Prepared Query exists", func(t *testing.T) {
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
			log.Printf("Creating NewSnippetModel failed")
			return
		}
		repo.LatestStatement = nil
		output, err := repo.Get(1)
		prep.ExpectQuery().WillReturnError(err)
		assert.Nil(t, output)
		assert.Error(t, err)
	})
}
