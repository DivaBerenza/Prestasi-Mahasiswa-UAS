package repository

import (
	"context"
	"fmt"
	"time"

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

	// Convert InsertedID ke ObjectID
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

// Update achievement by ID
func (r *AchievementRepository) Update(id string, updateData bson.M) (*model.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	updateData["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Achievement
	err = r.Collection.FindOneAndUpdate(r.Ctx, bson.M{"_id": oid}, bson.M{"$set": updateData}, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update achievement: %v", err)
	}

	return &updated, nil
}

// Delete achievement by ID (soft delete: update status to 'deleted')
func (r *AchievementRepository) Delete(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	_, err = r.Collection.UpdateOne(r.Ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"status":    "deleted",
		"updatedAt": time.Now(),
	}})
	if err != nil {
		return fmt.Errorf("failed to delete achievement: %v", err)
	}
	return nil
}

// UpdateStatus mengubah status achievement di MongoDB
func (r *AchievementRepository) UpdateStatus(id string, status string) (*model.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	update := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Achievement
	err = r.Collection.FindOneAndUpdate(r.Ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update status: %v", err)
	}

	return &updated, nil
}

func (r *AchievementRepository) GetVerifiedOnly() ([]model.Achievement, error) {
	filter := bson.M{"status": string(model.StatusVerified)}

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

// repository/achievement_repository.go
func (r *AchievementRepository) AddAttachment(
	achievementID string,
	attachment model.Attachment,
) error {

	objID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(
		r.Ctx,
		bson.M{"_id": objID},
		bson.M{
			"$push": bson.M{
				"attachments": attachment,
			},
		},
	)

	return err
}

func (r *AchievementRepository) GetStatistics(
	ctx context.Context,
) (
	total int64,
	submitted int64,
	verified int64,
	rejected int64,
	err error,
) {

	// total achievements
	total, err = r.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return
	}

	// per status
	submitted, err = r.Collection.CountDocuments(ctx, bson.M{"status": "submitted"})
	if err != nil {
		return
	}

	verified, err = r.Collection.CountDocuments(ctx, bson.M{"status": "verified"})
	if err != nil {
		return
	}

	rejected, err = r.Collection.CountDocuments(ctx, bson.M{"status": "rejected"})
	if err != nil {
		return
	}

	return
}

func (r *AchievementRepository) GetStatisticsByStudentID(
	ctx context.Context,
	studentID string,
) (
	total int64,
	submitted int64,
	verified int64,
	rejected int64,
	err error,
) {

	filter := bson.M{"studentId": studentID}

	total, err = r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	submitted, err = r.Collection.CountDocuments(ctx, bson.M{
		"studentId": studentID,
		"status":    "submitted",
	})
	if err != nil {
		return
	}

	verified, err = r.Collection.CountDocuments(ctx, bson.M{
		"studentId": studentID,
		"status":    "verified",
	})
	if err != nil {
		return
	}

	rejected, err = r.Collection.CountDocuments(ctx, bson.M{
		"studentId": studentID,
		"status":    "rejected",
	})
	if err != nil {
		return
	}

	return
}

