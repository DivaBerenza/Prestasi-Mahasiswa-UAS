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


func (r *AchievementReferenceRepository) StatusHistory(mongoID string) (*model.AchievementStatusHistory, error) {
	query := `
		SELECT status, updated_at, verified_by, rejection_note
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`

	var h model.AchievementStatusHistory
	err := r.DB.QueryRow(query, mongoID).Scan(
		&h.Status,
		&h.UpdatedAt,
		&h.VerifiedBy,
		&h.RejectionNote,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &h, nil
}

func (r *AchievementReferenceRepository) Submit(mongoID string) error {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET
			status = $1,
			submitted_at = $2,
			updated_at = $2
		WHERE mongo_achievement_id = $3
	`

	_, err := r.DB.Exec(
		query,
		model.StatusSubmitted,
		now,
		mongoID,
	)

	return err
}

func (r *AchievementReferenceRepository) Verify(
	mongoID string,
	verifiedBy uuid.UUID,
) error {

	now := time.Now()

	query := `
		UPDATE achievement_references
		SET
			status = $1,
			verified_at = $2,
			verified_by = $3,
			updated_at = $2
		WHERE mongo_achievement_id = $4
	`

	_, err := r.DB.Exec(
		query,
		model.StatusVerified,
		now,
		verifiedBy,
		mongoID,
	)

	return err
}

func (r *AchievementReferenceRepository) Reject(
	mongoID string,
	rejectionNote string,
) error {

	now := time.Now()

	query := `
		UPDATE achievement_references
		SET
			status = $1,
			verified_at = $2,
			rejection_note = $3,
			updated_at = $2
		WHERE mongo_achievement_id = $4
	`

	_, err := r.DB.Exec(
		query,
		model.StatusRejected,
		now,
		rejectionNote,
		mongoID,
	)

	return err
}





