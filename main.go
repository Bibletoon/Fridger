package main

import (
	"Fridger/internal/configuration"
	"Fridger/internal/handlers"
	"Fridger/internal/infrastructure/clients"
	"Fridger/internal/infrastructure/db"
	"Fridger/internal/services"
	"context"
	configuration_yaml_file "github.com/BoRuDar/configuration-yaml-file"
	configlib "github.com/BoRuDar/configuration/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "image/jpeg"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := configuration.Configuration{}
	configurator := configlib.New(
		&cfg,
		configlib.NewEnvProvider(),
		configuration_yaml_file.NewYAMLFileProvider("config.yaml"),
		configlib.NewDefaultProvider(),
	)
	err := configurator.InitValues()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	if cfg.MigrationsConfiguration.RunOnStart {
		err = db.Migrate(pool, cfg.MigrationsConfiguration)
		if err != nil {
			panic(err)
		}
	}

	crptClient := clients.NewCrptClient()
	photoHandler := handlers.NewPhotoHandler(crptClient)

	bot, err := services.NewBot(
		cfg,
		photoHandler)
	if err != nil {
		panic(err)
	}

	bot.Start(ctx)
}
