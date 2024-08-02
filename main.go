package main

import (
	"Fridger/internal/configuration"
	"Fridger/internal/handlers"
	"Fridger/internal/infrastructure/clients"
	"Fridger/internal/services"
	"context"
	configuration_yaml_file "github.com/BoRuDar/configuration-yaml-file"
	configlib "github.com/BoRuDar/configuration/v4"
	"github.com/go-telegram/bot"
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
		configuration_yaml_file.NewYAMLFileProvider("secrets.yaml"))
	err := configurator.InitValues()
	if err != nil {
		panic(err)
	}

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	crptClient := clients.NewCrptClient()
	photoHandler := handlers.NewPhotoHandler(crptClient)

	botService := services.NewBotService(
		b,
		photoHandler)
	botService.Start(ctx)
}
