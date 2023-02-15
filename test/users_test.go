package snippetbox_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"snippetbox/pkg/models/mysql"
	"testing"
)

func TestUsers(t *testing.T) {
	db, mock := NewMock()
	userModel := &mysql.UserModel{DB: db}
	t.Run("UserModel OK Case - Testing Insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users ...").WithArgs(
			"Name", "Email", sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := userModel.Insert("Name", "Email", "Password")
		assert.NoError(t, err)
	})
	t.Run("UserModel OK Case - Valid Credentials", func(t *testing.T) {
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
	t.Run("UserModel NOK Case - Special Error when looking for email", func(t *testing.T) {
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
	t.Run("UserModel NOK Case - No email found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "hashed_password"})
		rows.AddRow(
			1,
			"Password")

		mock.ExpectQuery(
			"SELECT id, hashed_password FROM users WHERE email \\= \\?").WithArgs(
			"NoEmail").WillReturnError(sql.ErrNoRows)
		authStatus, err := userModel.Authenticate("NoEmail", "Password")
		assert.Error(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("UserModel NOK Case - Invalid Password", func(t *testing.T) {
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

		userModel, err := userModel.Get(0)
		assert.NoError(t, err)
		assert.Nil(t, userModel)
	})
}
