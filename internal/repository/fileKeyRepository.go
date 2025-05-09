package repository

import (
	"aggron/internal/db/models"
	"aggron/internal/types"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FileKeyRepository struct {
	collection *mongo.Collection
}

func NewFileKeyRepository(db *mongo.Database) *FileKeyRepository {
	return &FileKeyRepository{
		collection: db.Collection("file_keys"),
	}
}

func (r *FileKeyRepository) CreateIndexes(ctx context.Context) error {
	// set ttl
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(types.DefaultExpirationTime.Seconds())),
	}

	_, err := r.collection.Indexes().CreateOne(ctx, ttlIndex)

	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func (r *FileKeyRepository) CreateFileKey(ctx context.Context, fileKey models.FileKey) (*models.FileKey, error) {
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
