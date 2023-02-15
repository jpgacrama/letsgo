package snippetbox_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/models/mysql"
	"testing"
)

func TestUsers(t *testing.T) {
	db, mock := NewMock()
	userModel := &mysql.UserModel{DB: db}
	t.Run("Insert OK Case - Testing Insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users ...").WithArgs(
			"Name", "Email", sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := userModel.Insert("Name", "Email", "Password")
		assert.NoError(t, err)
	})
	t.Run("Authenticate OK Case", func(t *testing.T) {
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
	t.Run("Authenticate NOK Case - Invalid Password", func(t *testing.T) {
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
