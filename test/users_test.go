package test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"
	"testing"
	"time"
)

func TestUserModelInsert(t *testing.T) {
	tests := []struct {
		testName       string
		userName       string
		email          string
		password       string
		duplicateEntry bool
	}{
		{
			testName:       "UserModel Insert OK Case",
			userName:       "Name",
			email:          "Email",
			password:       "Password",
			duplicateEntry: false,
		},
		{
			testName:       "UserModel Insert OK Case - Duplicate Entry",
			userName:       "Name",
			email:          "Email",
			password:       "Password",
			duplicateEntry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock := NewMock()
			userModel := &mysql.UserModel{DB: db}

			if tt.duplicateEntry {
				mysqlError := sqlDriver.MySQLError{
					Number:  1062,
					Message: "Error 1062: Duplicate entry",
				}
				mock.ExpectExec("INSERT INTO users ...").WithArgs(
					tt.userName, tt.email, sqlmock.AnyArg(),
				).WillReturnError(errors.New(mysqlError.Error())).WillReturnResult(driver.ResultNoRows)
			} else {
				mock.ExpectExec("INSERT INTO users ...").WithArgs(
					tt.userName, tt.email, sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := userModel.Insert(tt.userName, tt.email, tt.password)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsers(t *testing.T) {
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
