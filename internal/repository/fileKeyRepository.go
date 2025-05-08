package repository

import (
	"aggron/internal/db/models"
	"aggron/internal/types"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileKeyRepository struct {
	collection *mongo.Collection
}

func NewFileKeyRepository(db *mongo.Database) *FileKeyRepository {
	return &FileKeyRepository{
		collection: db.Collection("file_keys"),
	}
}

func (r *FileKeyRepository) CreateFileKey(ctx context.Context, fileKey models.FileKey) (*models.FileKey, error) {
	fileKey.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(types.DefaultExpirationTime))
	
	_, err := r.collection.InsertOne(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to insert file key: %w", err)
	}

	return &fileKey, nil
}

func (r *FileKeyRepository) FindByFileID(ctx context.Context, fileID string) (*models.FileKey, error) {
	var fileKey models.FileKey
	filter := bson.M{"file_id": fileID}

	err := r.collection.FindOne(ctx, filter).Decode(&fileKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding file key: %w", err)
	}

	return &fileKey, nil
}
