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
	StudentID          uuid.UUID         `json:"student_id"`
	MongoAchievementID string            `json:"mongo_achievement_id"`
	Status             AchievementStatus `json:"status"`
	SubmittedAt        *time.Time        `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time        `json:"verified_at,omitempty"`
	VerifiedBy         *uuid.UUID        `json:"verified_by,omitempty"`
	RejectionNote      *string           `json:"rejection_note,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}
