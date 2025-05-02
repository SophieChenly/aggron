package models

type User struct {
	DiscordID string `bson:"discord_id"`
	Email     string `bson:"email"`
}
