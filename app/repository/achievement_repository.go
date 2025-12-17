package repository

import (
	"context"

	"UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementRepository struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

func NewAchievementRepository(coll *mongo.Collection) *AchievementRepository {
	return &AchievementRepository{
		Collection: coll,
		Ctx:        context.Background(),
	}
}

// GetAll ambil semua achievements
func (r *AchievementRepository) GetAll() ([]model.Achievement, error) {
	cur, err := r.Collection.Find(r.Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(r.Ctx)

	var achievements []model.Achievement
	for cur.Next(r.Ctx) {
		var a model.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

// GetByStudentIDs ambil achievements berdasarkan student IDs
func (r *AchievementRepository) GetByStudentID(studentIDs []string) ([]model.Achievement, error) {
	filter := bson.M{"studentId": bson.M{"$in": studentIDs}}
	cur, err := r.Collection.Find(r.Ctx, filter, &options.FindOptions{
		Sort: bson.M{"createdAt": -1},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(r.Ctx)

	var achievements []model.Achievement
	for cur.Next(r.Ctx) {
		var a model.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}
