package services

import (
	"Fridger/internal/configuration"
	"Fridger/internal/domain/interfaces/handlers"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NewBot(cfg configuration.Configuration, handlers ...handlers.Handler) (*bot.Bot, error) {
	opt := []bot.Option{
		bot.WithMiddlewares(startGoroutineMiddleware),
	}

	b, err := bot.New(cfg.BotToken, opt...)

	if err != nil {
		return nil, err
	}

	for _, handler := range handlers {
		b.RegisterHandlerMatchFunc(handler.Match, handler.Handle)
	}

	return b, nil
}

func startGoroutineMiddleware(h bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		go h(ctx, bot, update)
	}
}
