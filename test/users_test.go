package snippetbox_test

import (
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
	t.Run("UserModel OK Case - Testing Authenticate", func(t *testing.T) {

		authStatus, err := userModel.Authenticate("Email", "Password")
		assert.NoError(t, err)
		assert.Equal(t, 0, authStatus)
	})
	t.Run("UserModel OK Case - Testing Get", func(t *testing.T) {

		userModel, err := userModel.Get(0)
		assert.NoError(t, err)
		assert.Nil(t, userModel)
	})
}
