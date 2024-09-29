package jobs

import (
	"Fridger/internal/configuration"
	"Fridger/internal/domain/interfaces/services"
	"Fridger/internal/helpers"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

type ExpiringProductsJob struct {
	cfg     configuration.ExpiringProductsJobConfiguration
	appCfg  configuration.AppConfiguration
	bot     *bot.Bot
	service services.ProductService
}

func NewExpiringProductsJob(
	cfg configuration.ExpiringProductsJobConfiguration,
	appCfg configuration.AppConfiguration,
	bot *bot.Bot,
	service services.ProductService) *ExpiringProductsJob {
	return &ExpiringProductsJob{
		cfg:     cfg,
		appCfg:  appCfg,
		bot:     bot,
		service: service,
	}
}

func (j *ExpiringProductsJob) Run() {
	ctx := context.Background()

	daysBeforeExpiration := j.appCfg.DaysBeforeExpiration
	products, err := j.service.GetExpiringProducts(ctx, daysBeforeExpiration)

	if err != nil {
		slog.Error("Error getting expiring products", "error", err.Error())
		return
	}

	msgText := helpers.BuildExpiredListMessage(products)

	if len(msgText) == 0 {
		return
	}

	msg := &bot.SendMessageParams{
		ChatID:    j.cfg.NotificationUserId,
		Text:      msgText,
		ParseMode: models.ParseModeHTML,
	}

	_, err = j.bot.SendMessage(ctx, msg)

	if err != nil {
		slog.Error("Error sending expiring products", "error", err.Error())
	}
}
