package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"studentId"`
	AchievementType string             `bson:"achievementType" json:"achievementType"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`

	Details bson.M `bson:"details,omitempty" json:"details"`

	Attachments []Attachment `bson:"attachments,omitempty" json:"attachments"`
	Tags        []string     `bson:"tags,omitempty" json:"tags"`
	Points      int          `bson:"points" json:"points"`

	Status    string    `bson:"status" json:"status"` // draft, submitted, verified, rejected
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FilePath    string    `bson:"filePath" json:"filePath"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}
