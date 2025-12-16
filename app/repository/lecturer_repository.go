package repository

import (
	"database/sql"
	"UAS/app/model"
)

type LecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepository(db *sql.DB) *LecturerRepository {
	return &LecturerRepository{DB: db}
}

func (r *LecturerRepository) GetAll() ([]model.Lecturer, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		lecturers = append(lecturers, l)
	}

	return lecturers, nil
}
