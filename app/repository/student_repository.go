package repository

import (
	"UAS/app/model"
	"database/sql"
	"fmt"
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

// GetStudentByID mengambil satu mahasiswa berdasarkan UUID student
func (r *StudentRepository) GetByID(studentID string) (*model.StudentResponse, error) {
	row := r.DB.QueryRow(`
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1
	`, studentID)

	var s model.StudentResponse
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.FullName,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, err
	}

	return &s, nil
}

// UpdateAdvisor assign atau update advisor_id mahasiswa
func (r *StudentRepository) UpdateAdvisor(studentID string, advisorID string) (*model.StudentResponse, error) {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id
	`

	row := r.DB.QueryRow(query, advisorID, studentID)

	var s model.StudentResponse
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, err
	}

	// Ambil full_name dari tabel users
	err = r.DB.QueryRow(`SELECT full_name FROM users WHERE id = $1`, s.UserID).Scan(&s.FullName)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetByUserID mengambil mahasiswa berdasarkan user_id
func (r *StudentRepository) GetByUserID(userID string) ([]model.StudentResponse, error) {
	rows, err := r.DB.Query(`
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.user_id = $1
	`, userID)
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

// GetByAdvisorID mengambil mahasiswa berdasarkan advisor_id
func (r *StudentRepository) GetByAdvisorID(advisorID string) ([]model.StudentResponse, error) {
	rows, err := r.DB.Query(`
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
	`, advisorID)
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


