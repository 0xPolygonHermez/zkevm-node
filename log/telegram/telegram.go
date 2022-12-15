package telegram

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramBot is a struct for the telegram communication.
type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	mu     sync.Mutex
}

// NewBot creates a new telegram bot.
func NewBot(cfg Config) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.APIToken)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		bot:    bot,
		chatID: cfg.ChatID,
	}, err
}

func (t *TelegramBot) sendMessage(level string, message string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	msg := tgbotapi.NewMessage(t.chatID, fmt.Sprintf("*%s:* `%s`", level, message))
	msg.ParseMode = "Markdown"
	_, err := t.bot.Send(msg)
	return err
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

// Info sends a info level message.
func (t *TelegramBot) Info(args ...interface{}) error {
	return t.sendMessage("INFO", getMessage("", args))
}

// Warn sends a warnning level message.
func (t *TelegramBot) Warn(args ...interface{}) error {
	return t.sendMessage("WARN", getMessage("", args))
}

// Error sends an error level message.
func (t *TelegramBot) Error(args ...interface{}) error {
	return t.sendMessage("ERROR", getMessage("", args))
}

// Debug sends a debug level message.
func (t *TelegramBot) Debug(args ...interface{}) error {
	return t.sendMessage("DEBUG", getMessage("", args))
}

// Fatal sends a fatal level message.
func (t *TelegramBot) Fatal(args ...interface{}) error {
	return t.sendMessage("FATAL", getMessage("", args))
}

// Infof sends log.Infof.
func (t *TelegramBot) Infof(template string, args ...interface{}) error {
	return t.Info(getMessage(template, args))
}

// Warnf sends log.Warnf.
func (t *TelegramBot) Warnf(template string, args ...interface{}) error {
	return t.Warn(getMessage(template, args))
}

// Errorf sends log.Errorf.
func (t *TelegramBot) Errorf(template string, args ...interface{}) error {
	return t.Error(getMessage(template, args))
}

// Debugf sends log.Debugf.
func (t *TelegramBot) Debugf(template string, args ...interface{}) error {
	return t.Debug(getMessage(template, args))
}

// Fatalf sends log.Debugf.
func (t *TelegramBot) Fatalf(template string, args ...interface{}) error {
	return t.Fatal(getMessage(template, args))
}
