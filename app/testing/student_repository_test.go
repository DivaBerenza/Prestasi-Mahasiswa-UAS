package repository_test

import (
	"testing"

	"UAS/app/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllStudents_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)
	studentID := uuid.New()
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "full_name", "program_study", "academic_year", "advisor_id"}).
		AddRow(studentID, userID, "20221001", "Student One", "TI", "2022", nil)

	mock.ExpectQuery("SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id FROM students s JOIN users u ON u.id = s.user_id ORDER BY s.created_at ASC").
		WillReturnRows(rows)

	students, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Student One", students[0].FullName)
}

func TestGetByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)
	studentID := uuid.New()
	userID := uuid.New()

	// Mock student row
	mockStudentRow := sqlmock.NewRows([]string{"id", "user_id", "student_id", "full_name", "program_study", "academic_year", "advisor_id"}).
		AddRow(studentID, userID, "20221001", "Student One", "TI", "2022", nil)

	mock.ExpectQuery("SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, s.advisor_id FROM students s JOIN users u ON u.id = s.user_id WHERE s.id = \\$1").
		WithArgs(studentID.String()).
		WillReturnRows(mockStudentRow)

	student, err := repo.GetByID(studentID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Student One", student.FullName)
}
