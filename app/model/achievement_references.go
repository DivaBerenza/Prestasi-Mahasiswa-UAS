package model

import (
	"time"

	"github.com/google/uuid"
)

type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
)

type AchievementReference struct {
	ID                 uuid.UUID         `json:"id"`
	StudentID          uuid.UUID         `json:"student_id"`             // FK -> students.id
	MongoAchievementID string            `json:"mongo_achievement_id"`    // Document ID in Mongo (24 chars)
	Status             AchievementStatus `json:"status"`                  // ENUM value

	SubmittedAt *time.Time `json:"submitted_at,omitempty"` // nullable
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`  // nullable
	VerifiedBy  *uuid.UUID `json:"verified_by,omitempty"`  // FK -> users.id, nullable
	RejectionNote *string  `json:"rejection_note,omitempty"` // null if not rejected

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
