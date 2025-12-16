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

func (r *LecturerRepository) GetAll() ([]model.LecturerResponse, error) {
	rows, err := r.DB.Query(`
		SELECT l.id, l.user_id, u.full_name, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY l.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.LecturerResponse
	for rows.Next() {
		var l model.LecturerResponse
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.FullName,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, l)
	}

	return result, nil
}

func (r *LecturerRepository) GetAdvisees(lecturerID string) ([]model.AdviseeResponse, error) {
	rows, err := r.DB.Query(`
		SELECT s.student_id, u.full_name, s.program_study, s.academic_year
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advisees []model.AdviseeResponse
	for rows.Next() {
		var a model.AdviseeResponse
		if err := rows.Scan(
			&a.StudentID,
			&a.FullName,
			&a.ProgramStudy,
			&a.AcademicYear,
		); err != nil {
			return nil, err
		}
		advisees = append(advisees, a)
	}

	return advisees, nil
}