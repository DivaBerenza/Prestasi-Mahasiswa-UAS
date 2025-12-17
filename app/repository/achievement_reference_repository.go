package repository

import (
	"database/sql"
	"fmt"
	"time"

	"UAS/app/model"

	"github.com/google/uuid"
)

type AchievementReferenceRepository struct {
	DB *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{DB: db}
}

// CreateReference membuat reference prestasi baru
func (r *AchievementReferenceRepository) CreateReference(studentUUID uuid.UUID, mongoID string) (*model.AchievementReference, error) {
	query := `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	`
	id := uuid.New()
	now := time.Now()
	ref := &model.AchievementReference{}

	err := r.DB.QueryRow(query, id, studentUUID, mongoID, "draft", now, now).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed create reference: %v", err)
	}

	return ref, nil
}

