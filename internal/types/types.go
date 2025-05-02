package types

type StateInfo struct {
	SenderDiscordID   string
	SenderEmail       string
	ReceiverDiscordID string
	ReceiverEmail     string
	FileID            string
	State             string
}

type AuthInfo struct {
	DiscordID string
	Email     string
}
