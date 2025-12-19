package repository_test

import (
	"testing"
	"UAS/app/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewUserRepository(db)
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "full_name", "role_id", "is_active"}).
		AddRow(userID, "user1", "user1@example.com", "hashedpass", "User One", uuid.New(), true)

	mock.ExpectQuery("SELECT id, username, email, password, full_name, role_id, is_active FROM users WHERE email = \\$1").
		WithArgs("user1@example.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("user1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "User One", user.FullName)
}
