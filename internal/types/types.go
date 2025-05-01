package types

type StateInfo struct {
	SenderDiscordID   string
	ReceiverDiscordID string
	FileID            string
	State             string
}

type AuthInfo struct {
	DiscordID string
	State     string
}
