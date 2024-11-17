package services

import (
	"Fridger/internal/configuration"
	"Fridger/internal/domain/interfaces/handlers"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"slices"
)

func NewBot(cfg configuration.BotConfiguration, handlers ...handlers.Handler) (*bot.Bot, error) {
	userFilterMiddleware := createUserFilterMiddleware(cfg)

	opt := []bot.Option{
		bot.WithMiddlewares(userFilterMiddleware, startGoroutineMiddleware),
	}

	b, err := bot.New(cfg.Token, opt...)

	if err != nil {
		return nil, err
	}

	for _, handler := range handlers {
		b.RegisterHandlerMatchFunc(handler.Match, handler.Handle)
	}

	return b, nil
}

func createUserFilterMiddleware(cfg configuration.BotConfiguration) bot.Middleware {
	return func(h bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if !slices.Contains(cfg.AllowedUserIds, update.Message.From.ID) {
				return
			}

			h(ctx, bot, update)
		}
	}
}

func startGoroutineMiddleware(h bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		go h(ctx, bot, update)
	}
}
