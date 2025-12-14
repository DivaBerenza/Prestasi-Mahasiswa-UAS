package repository

import (
	"database/sql"
	"errors"

	"UAS/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Cari user berdasarkan email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, email, password, full_name, role_id, is_active
		FROM users
		WHERE email = $1
	`, email)

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

// Ambil permissions user berdasarkan userID
func (r *UserRepository) GetPermissionsByUserID(userID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT p.resource || ':' || p.action
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}

	return perms, nil
}

// repository/user_repository.go
func (r *UserRepository) GetAllUsers() ([]*model.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, username, email, full_name, role_id, is_active
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}


