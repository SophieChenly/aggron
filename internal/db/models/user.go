package models

type User struct {
	DiscordID string `bson:"discord_id"`
	PassageID string `bson:"passage_id"`
}
