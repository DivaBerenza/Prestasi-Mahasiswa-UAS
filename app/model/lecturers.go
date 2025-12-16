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


