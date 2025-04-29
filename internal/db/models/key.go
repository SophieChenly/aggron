package models

type Key struct {
	DiscordID         string `bson:"discord_id"`
	EncryptedKey      string `bson:"encrypted_key"`
	FileID            string `bson:"file_id"`
	ReceiverDiscordID string `bson:"receiver_discord_id"`
}
