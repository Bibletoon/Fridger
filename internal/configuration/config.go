package configuration

type Configuration struct {
	BotConfiguration        BotConfiguration        `file_yml:"bot"`
	ConnectionString        string                  `file_yml:"connection_string" env:"CONNECTION_STRING"`
	MigrationsConfiguration MigrationsConfiguration `file_yml:"migrations_configuration"`
	JobsConfiguration       JobsConfiguration       `file_yml:"jobs"`
}

type BotConfiguration struct {
	Token          string  `file_yml:"bot.token" env:"BOT_TOKEN"`
	AllowedUserIds []int64 `file_yml:"bot.allowed_user_ids" env:"BOT_USER_ID"`
}

type MigrationsConfiguration struct {
	RunOnStart       bool   `file_yml:"migrations_configuration.run_on_start" default:"true"`
	MigrationsFolder string `file_yml:"migrations_configuration.migrations_folder" default:"migrations"`
}

type JobsConfiguration struct {
	ExpiringProductsJobConfiguration ExpiringProductsJobConfiguration `file_yml:"jobs.expiring_products"`
}

type ExpiringProductsJobConfiguration struct {
	Schedule             string `file_yml:"jobs.expiring_products.schedule" default:"0 5,10,14 * * *"`
	DaysBeforeExpiration int    `file_yml:"jobs.expiring_products.days_before_expiration" default:"3"`
	NotificationUserId   int64  `file_yml:"jobs.expiring_products.notification_user_id"`
}
