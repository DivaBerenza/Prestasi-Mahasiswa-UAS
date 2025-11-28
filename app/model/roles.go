package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"role"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

