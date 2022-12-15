package telegram

// Config for telegram bot
type Config struct {
	// APIToken is a token for the telegram api access.
	APIToken string
	// ChatID indicates the global id of the chat bot.
	ChatID int64
}
