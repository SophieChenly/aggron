package types

import "time"

type StateInfo struct {
	SenderDiscordID   string
	SenderEmail       string
	ReceiverDiscordID string
	ReceiverEmail     string
	FileID            string
}

type AuthInfo struct {
	DiscordID string
	Email     string
}

var DefaultExpirationTime time.Duration = time.Minute * 15
