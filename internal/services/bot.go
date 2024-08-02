package services

import (
	"Fridger/internal/domain/interfaces/handlers"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotService struct {
	bot *bot.Bot
}

func NewBotService(b *bot.Bot, handlers ...handlers.Handler) *BotService {
	wrapper := &BotService{
		bot: b,
	}

	for _, handler := range handlers {
		b.RegisterHandlerMatchFunc(handler.Match, func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			go handler.Handle(ctx, bot, update)
		})
	}

	return wrapper
}

func (b *BotService) Start(ctx context.Context) {
	b.bot.Start(ctx)
}
