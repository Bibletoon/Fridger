package jobs

import (
	"Fridger/internal/configuration"
	"Fridger/internal/domain/interfaces/services"
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log/slog"
)

type ExpiringProductsJob struct {
	cfg     configuration.ExpiringProductsJobConfiguration
	bot     *bot.Bot
	service services.ProductService
}

func NewExpiringProductsJob(cfg configuration.ExpiringProductsJobConfiguration, bot *bot.Bot, service services.ProductService) *ExpiringProductsJob {
	return &ExpiringProductsJob{
		cfg:     cfg,
		bot:     bot,
		service: service,
	}
}

func (j *ExpiringProductsJob) Run() {
	ctx := context.Background()

	daysBeforeExpiration := j.cfg.DaysBeforeExpiration
	expiringProducts, expiredProducts, err := j.service.GetExpiringProducts(ctx, daysBeforeExpiration)

	if err != nil {
		slog.Error("Error getting expiring products", "error", err.Error())
		return
	}

	if len(expiredProducts) == 0 {
		return
	}

	buf := bytes.Buffer{}

	if len(expiringProducts) > 0 {
		buf.WriteString("*Продукты, которые скоро испортятся:*\n\n")

		for i, product := range expiringProducts {
			buf.WriteString(fmt.Sprintf("%d. %s %s\n", i, product.Name, product.CreatedAt.String()))
		}
	}

	if len(expiredProducts) > 0 {
		buf.WriteString("*Продукты, которые уже испортились:*\n\n")

		for i, product := range expiringProducts {
			buf.WriteString(fmt.Sprintf("%d. %s %s", i, product.Name, product.CreatedAt.String()))
		}
	}

	params := &bot.SendMessageParams{
		ChatID: j.cfg.NotificationUserId,
		Text:   buf.String(),
	}

	_, err = j.bot.SendMessage(ctx, params)

	if err != nil {
		slog.Error("Error sending expiring products", "error", err.Error())
	}
}
