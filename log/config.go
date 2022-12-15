package log

import "github.com/0xPolygonHermez/zkevm-node/log/telegram"

// Config for log
type Config struct {
	// Environment defining the log format ("production" or "development").
	Environment LogEnvironment `mapstructure:"Environment"`
	// Level of log, e.g. INFO, WARN, ...
	Level string `mapstructure:"Level"`
	// Outputs
	Outputs []string `mapstructure:"Outputs"`
	// Receiver represents the kind of monitoring app ("telegram" or "slack").
	Receiver string
	// TelegramConfig is a config for the telegram bot.
	TelegramConfig telegram.Config
}
