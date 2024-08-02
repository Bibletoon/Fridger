package configuration

type Configuration struct {
	BotToken string `file_yml:"bot_token" env:"BOT_TOKEN"`
}
