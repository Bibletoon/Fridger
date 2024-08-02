package main

import (
	bot2 "Fridger/internal/bot"
	"Fridger/internal/infrastructure/clients"
	"context"
	"github.com/go-telegram/bot"
	_ "image/jpeg"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New("")
	if err != nil {
		panic(err)
	}

	crptClient := clients.NewCrptClient()
	botWrapper := bot2.NewBotWrapper(b, crptClient)
	botWrapper.Start(ctx)
}
