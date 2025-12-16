package repository

import (
	"UAS/app/model"
	"database/sql"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

// GetAllStudents mengambil semua mahasiswa
func (r *StudentRepository) GetAll() ([]model.StudentResponse, error) {
	rows, err := r.DB.Query(`
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id
		FROM students s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.StudentResponse
	for rows.Next() {
		var s model.StudentResponse
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.FullName,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}
