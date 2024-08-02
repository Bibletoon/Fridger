package bot

import (
	"Fridger/internal/domain/interfaces/clients"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotWrapper struct {
	bot        *bot.Bot
	crptClient clients.CrptClient
}

func NewBotWrapper(bot *bot.Bot, crptClient clients.CrptClient) *BotWrapper {
	wrapper := &BotWrapper{
		bot:        bot,
		crptClient: crptClient,
	}
	registerHandlers(wrapper)
	return wrapper
}

func registerHandlers(w *BotWrapper) {
	w.bot.RegisterHandlerMatchFunc(
		handlePhotoMatchFunc,
		func(ctx context.Context, bot *bot.Bot, upd *models.Update) {
			go w.handlePhoto(ctx, upd)
		})
}

func (b *BotWrapper) Start(ctx context.Context) {
	b.bot.Start(ctx)
}
