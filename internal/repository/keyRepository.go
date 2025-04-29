package repository

import (
	"aggron/internal/db/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type KeyRepository struct {
	collection *mongo.Collection
}

func NewKeyRepository(db *mongo.Database) *KeyRepository {
	return &KeyRepository{
		collection: db.Collection("keys"),
	}
}

func (r *KeyRepository) CreateKey(ctx context.Context, discordID, encryptedKey, fileID, receiverDiscordID string) (*models.Key, error) {
	Key := models.Key{
		DiscordID:         discordID,
		EncryptedKey:      encryptedKey,
		FileID:            fileID,
		ReceiverDiscordID: receiverDiscordID,
	}

	_, err := r.collection.InsertOne(ctx, Key)
	if err != nil {
		return nil, err
	}

	return &Key, nil
}

func (r *KeyRepository) FindByDiscordIDAndFileID(ctx context.Context, discordID, fileID string) (*models.Key, error) {
	var Key models.Key
	filter := bson.M{"discord_id": discordID, "file_id": fileID}

	err := r.collection.FindOne(ctx, filter).Decode(&Key)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return &Key, nil
}
