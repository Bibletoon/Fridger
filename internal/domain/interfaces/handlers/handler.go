package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler interface {
	Match(update *models.Update) bool
	Handle(ctx context.Context, bot *bot.Bot, update *models.Update)
}
