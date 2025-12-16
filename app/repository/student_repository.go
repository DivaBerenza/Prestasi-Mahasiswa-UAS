package repository

import (
	"database/sql"
	"UAS/app/model"

	"github.com/google/uuid"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

// Get all students
func (r *StudentRepository) GetAll() ([]model.Student, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []model.Student{}

	for rows.Next() {
		var s model.Student
		var idStr, userIDStr, advisorIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisorIDStr, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.ID, _ = uuid.Parse(idStr)
		s.UserID, _ = uuid.Parse(userIDStr)
		s.AdvisorID, _ = uuid.Parse(advisorIDStr)
		students = append(students, s)
	}

	return students, nil
}

// GetByUserID ambil student berdasarkan user_id
func (r *StudentRepository) GetByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT id, student_id
		FROM students
		WHERE user_id = $1
	`
	row := r.DB.QueryRow(query, userID)

	var s model.Student
	var idStr string
	if err := row.Scan(&idStr, &s.StudentID); err != nil {
		return nil, err
	}

	s.ID, _ = uuid.Parse(idStr)
	return &s, nil
}

func (r *StudentRepository) CreateStudent(student *model.Student) error {
	_, err := r.DB.Exec(`
		INSERT INTO students (
			id, user_id, student_id, program_study, academic_year, advisor_id
		)
		VALUES ($1,$2,$3,$4,$5,$6)
	`,
		student.ID,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID, // boleh NULL
	)
	return err
}



