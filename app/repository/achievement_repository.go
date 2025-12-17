package repository

import (
	"context"
	"fmt"

	"UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *AchievementRepository) Create(achievement *model.Achievement) (*model.Achievement, error) {
	result, err := r.Collection.InsertOne(r.Ctx, achievement)
	if err != nil {
		return nil, err
	}

	// Aman: convert InsertedID ke primitive.ObjectID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert insertedID to ObjectID")
	}
	achievement.ID = oid

	return achievement, nil
}

func (r *AchievementRepository) GetByID(id string) (*model.Achievement, error) {
	// Convert string ID ke ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var achievement model.Achievement
	err = r.Collection.FindOne(r.Ctx, bson.M{"_id": oid}).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // tidak ditemukan
		}
		return nil, fmt.Errorf("failed to fetch achievement: %v", err)
	}

	return &achievement, nil
}

