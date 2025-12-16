package model

import (
	"time"

	"github.com/google/uuid"
)

type Lecturer struct {
	ID         uuid.UUID `json:"id"`
	UserID	 uuid.UUID `json:"user_id"`
	LecturerID   string    `json:"lecturer_id"`    // example: "Dosen001"
	Department   string    `json:"department"`     // example: "Teknik Informatika"
	CreatedAt  time.Time `json:"created_at"`
}

type LecturerResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	FullName   string    `json:"full_name"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdviseeResponse struct {
	StudentID    string `json:"student_id"`
	FullName     string `json:"full_name"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
}

