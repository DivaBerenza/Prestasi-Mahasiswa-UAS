package repository

import (
	"context"
	"time"


	"UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive" 
)

type AchievementRepository struct {
	Collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		Collection: db.Collection("achievements"),
	}
}

// Get achievements filtered by role and optional studentId
func (r *AchievementRepository) GetAchievements(role, studentId, status, achType string, limit, page int64) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if role == "student" {
		if studentId != "" {
			filter["studentId"] = studentId
		}
	} else if role == "admin" {
		if studentId != "" {
			filter["studentId"] = studentId
		}
		if status != "" {
			filter["status"] = status
		}
		if achType != "" {
			filter["achievementType"] = achType
		}
	}

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"createdAt": -1})

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementRepository) GetAchievementByID(id string) (*model.Achievement, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var achievement model.Achievement
    if err := r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&achievement); err != nil {
        return nil, err
    }

    return &achievement, nil
}

func (r *AchievementRepository) CreateAchievement(a *model.Achievement) (*model.Achievement, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    a.ID = primitive.NewObjectID()
    a.CreatedAt = time.Now()
    a.UpdatedAt = time.Now()
    a.Status = "draft" // status awal

    _, err := r.Collection.InsertOne(ctx, a)
    if err != nil {
        return nil, err
    }
    return a, nil
}



