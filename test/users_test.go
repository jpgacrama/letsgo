package snippetbox_test

import (
	"database/sql"
	"fmt"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/models"
)

func TestUsers(t *testing.T) {
	t.Run("Insert OK Case - Testing Insert", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		mock.ExpectExec("INSERT INTO users ...").WithArgs(
			"Name", "Email", sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := userModel.Insert("Name", "Email", "Password")
		assert.NoError(t, err)
	})
	t.Run("Authenticate OK Case", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		password := "C0mpl3xPass!"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "hashed_password"})
		rows.AddRow(
			1,
			hashedPassword)

		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"Email").WillReturnRows(rows)

		authStatus, err := userModel.Authenticate("Email", password)
		assert.NoError(t, err)
		assert.Equal(t, 1, authStatus)
	})
	t.Run("Authenticate NOK Case - Password is too short", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		password := "password123"
		rows := sqlmock.NewRows([]string{"id", "hashed_password"})
		rows.AddRow(
			1,
			password)

		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"Email").WillReturnRows(rows)
		authStatus, err := userModel.Authenticate("Email", password)
		assert.Error(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("Authenticate NOK Case - Special Error when looking for email", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		rows := sqlmock.NewRows([]string{"id", "hashed_password"})
		rows.AddRow(
			1,
			"Password")

		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"NoEmail").WillReturnRows(rows)
		authStatus, err := userModel.Authenticate("Email", "Password")
		assert.Error(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("Authenticate NOK Case - No email found", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"NoEmail").WillReturnError(sql.ErrNoRows)
		authStatus, err := userModel.Authenticate("NoEmail", "Password")
		assert.Error(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("Authenticate NOK Case - Invalid Password", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		rows := sqlmock.NewRows([]string{"id", "hashed_password"})
		rows.AddRow(
			1,
			"$2a$12$VRScRZpYfMbhPbF6qP/4le9kwuYO.VHPiugOtf62VKsITwRi2wGHS")

		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"Email").WillReturnRows(rows)
		authStatus, err := userModel.Authenticate("Email", "Password")
		assert.Error(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("UserModel OK Case - Testing Get", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		id := 1
		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "created"})
		timeCreated, err := time.Parse(time.RFC3339, "2024-02-23T10:23:42Z")
		if err != nil {
			fmt.Printf("parsing time failed")
		}

		rows.AddRow(
			id, "Jonas", "jonas@email.com", timeCreated)
		mock.ExpectQuery(
			"SELECT id, name, email, created FROM users WHERE id \\= \\?").
			WithArgs(id).WillReturnRows(rows)
		modelsUser, newErr := userModel.Get(id)
		assert.NoError(t, newErr)
		assert.NotNil(t, modelsUser)
	})
	t.Run("UserModel NOK Case - Testing Get - No Rows", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		id := 1

		mock.ExpectQuery(
			"SELECT id, name, email, created FROM users WHERE id \\= \\?").
			WithArgs(id).WillReturnError(sql.ErrNoRows)
		modelsUser, newErr := userModel.Get(id)
		assert.Error(t, newErr)
		assert.Nil(t, modelsUser)
	})
	t.Run("UserModel NOK Case - Testing Get - Unexpected Error", func(t *testing.T) {
		db, mock := NewMock()
		userModel := &mysql.UserModel{DB: db}
		id := 1

		mock.ExpectQuery(
			"SELECT id, name, email, created FROM users WHERE id \\= \\?").
			WithArgs(id).WillReturnError(models.ErrInvalidCredentials)
		modelsUser, newErr := userModel.Get(id)
		assert.Error(t, newErr)
		assert.Nil(t, modelsUser)
	})
}
