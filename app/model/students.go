package model

import (
	"time"

	"github.com/google/uuid"

)

type Student struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	StudentID    string    `json:"student_id"`     // example: "20221030001"
	ProgramStudy string    `json:"program_study"`  // example: "Teknik Informatika"
	AcademicYear string    `json:"academic_year"`  // example: "2025"
	AdvisorID    uuid.UUID `json:"advisor_id"`     // Lecturer reference

	CreatedAt time.Time `json:"created_at"`

}