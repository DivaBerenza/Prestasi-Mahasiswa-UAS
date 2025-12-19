package repository

import (
	"UAS/app/model"
	"database/sql"
)

type ReportRepository struct {
	DB *sql.DB
}

// Constructor
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

//
// ============================
// Statistik Mahasiswa (PostgreSQL ONLY)
// ============================
//
func (r *ReportRepository) GetStudentStatistics() (int, error) {
	var totalStudents int
	query := `SELECT COUNT(*) FROM students`

	err := r.DB.QueryRow(query).Scan(&totalStudents)
	if err != nil {
		return 0, err
	}

	return totalStudents, nil
}

//
// ============================
// Data Mahasiswa (untuk report detail)
// ============================
//
func (r *ReportRepository) GetStudentBase(studentID string) (*model.ReportStudent, error) {
	query := `
		SELECT
			s.student_id,
			u.username
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.student_id = $1
	`

	row := r.DB.QueryRow(query, studentID)

	var report model.ReportStudent
	err := row.Scan(&report.StudentID, &report.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *ReportRepository) GetStudentIDByNIM(nim string) (string, error) {
	query := `
		SELECT id
		FROM students
		WHERE student_id = $1
	`

	var studentUUID string
	err := r.DB.QueryRow(query, nim).Scan(&studentUUID)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return studentUUID, nil
}


