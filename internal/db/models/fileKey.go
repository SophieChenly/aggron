package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileKey struct {
	FileID            string           `bson:"file_id"`
	EncryptedKey      primitive.Binary `bson:"encrypted_key"`
	FileHash          primitive.Binary `bson:"file_hash"`
	SenderDiscordID   string           `bson:"sender_discord_id"`
	ReceiverDiscordID string           `bson:"receiver_discord_id"`
	CreatedAt         time.Time        `bson:"created_at"`
}
