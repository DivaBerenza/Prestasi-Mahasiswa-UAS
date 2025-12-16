package repository

import (
	"database/sql"
	"log"
	
	"github.com/google/uuid"
)

type AchievementReferenceRepository struct {
	DB *sql.DB
}

// Constructor
func NewAchievementReferenceRepository(db *sql.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{DB: db}
}

// Create reference
func (r *AchievementReferenceRepository) Create(studentID string, mongoID string) error {
	query := `
    INSERT INTO achievement_references
    (id, student_id, mongo_achievement_id, status, created_at, updated_at)
    VALUES ($1, $2, $3, 'draft', NOW(), NOW())
    `
	_, err := r.DB.Exec(query, uuid.New(), studentID, mongoID)
	if err != nil {
		log.Println("❌ Failed to insert achievement_reference:", err)
	}
	return err
}
