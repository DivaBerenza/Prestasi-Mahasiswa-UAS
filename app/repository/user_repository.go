package repository

import (
	"database/sql"
	"errors"
	// "time"
	"fmt"

	"UAS/app/model"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
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
		ORDER BY created_at ASC
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

// GetUserByID mengambil user berdasarkan UUID
func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, email, password, full_name, role_id, is_active
		FROM users
		WHERE id = $1
	`, id)

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

func (r *UserRepository) CreateUser(
	user *model.User,
	lecturerID *string,
	department *string,
	studentID *string,
	programStudy *string,
	academicYear *string,
) (*model.User, error) {

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user.ID = uuid.New()

	// ===== INSERT USER =====
	err = tx.QueryRow(`
		INSERT INTO users (
			id, username, email, password, full_name, role_id, is_active
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING created_at, updated_at
	`,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.FullName,
		user.RoleID,
		user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// ===== INSERT STUDENT =====
	if studentID != nil {
		_, err := tx.Exec(`
			INSERT INTO students (
				id, user_id, student_id, program_study, academic_year, advisor_id
			) VALUES ($1,$2,$3,$4,$5,NULL)
		`,
			uuid.New(),
			user.ID,
			*studentID,
			*programStudy,
			*academicYear,
		)
		if err != nil {
			return nil, err
		}
	}

	// ===== INSERT LECTURER (kalau ada) =====
	if lecturerID != nil {
		_, err := tx.Exec(`
			INSERT INTO lecturers (
				id, user_id, lecturer_id, department
			) VALUES ($1,$2,$3,$4)
		`,
			uuid.New(),
			user.ID,
			*lecturerID,
			*department,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}



func (r *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET username = $1,
		    email = $2,
		    full_name = $3,
		    role_id = $4,
		    is_active = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING id, username, email, full_name, role_id, is_active, created_at, updated_at
	`

	updatedUser := &model.User{}
	err := r.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.RoleID,
		&updatedUser.IsActive,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) (*model.User, error) {
	query := `
		UPDATE users
		SET password = $1,
		    updated_at = NOW()
		WHERE id = $2
		RETURNING id, username, email, full_name, role_id, is_active, created_at, updated_at
	`

	updatedUser := &model.User{}
	err := r.DB.QueryRow(query, hashedPassword, userID).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.RoleID,
		&updatedUser.IsActive,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) GetRoleIDByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.DB.QueryRow(`
		SELECT id FROM roles WHERE name = $1
	`, name).Scan(&id)
	return id, err
}






