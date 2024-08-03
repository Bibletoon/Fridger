package configuration

type Configuration struct {
	BotToken                string                  `file_yml:"bot_token" env:"BOT_TOKEN"`
	ConnectionString        string                  `file_yml:"connection_string" env:"CONNECTION_STRING"`
	MigrationsConfiguration MigrationsConfiguration `file_yml:"migrations_configuration"`
}

type MigrationsConfiguration struct {
	RunOnStart       bool   `file_yml:"migrations_configuration.run_on_start" default:"true"`
	MigrationsFolder string `file_yml:"migrations_configuration.migrations_folder" default:"migrations"`
}
