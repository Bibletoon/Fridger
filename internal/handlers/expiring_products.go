package handlers

import (
	"Fridger/internal/configuration"
	"Fridger/internal/domain/interfaces/handlers"
	"Fridger/internal/domain/interfaces/services"
	"Fridger/internal/helpers"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"strings"
)

type expiringProductsHandler struct {
	productService services.ProductService
	cfg            configuration.AppConfiguration
}

func NewExpiringProductsHandler(productService services.ProductService, cfg configuration.AppConfiguration) handlers.Handler {
	return &expiringProductsHandler{
		productService: productService,
		cfg:            cfg,
	}
}

func (h *expiringProductsHandler) Match(update *models.Update) bool {
	return strings.HasPrefix(update.Message.Text, "/expiring")
}

func (h *expiringProductsHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	daysBeforeExpiration := h.cfg.DaysBeforeExpiration
	products, err := h.productService.GetExpiringProducts(ctx, daysBeforeExpiration)

	if err != nil {
		slog.Error("Error getting expiring products", "error", err.Error())
		return
	}

	msgText := helpers.BuildExpiredProductsListMessage(products)

	if len(msgText) == 0 {
		msgText = "Продукты с опасным сроком годности не найдены"
	}

	msg := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msgText,
		ParseMode: models.ParseModeHTML,
	}

	_, err = b.SendMessage(ctx, msg)

	if err != nil {
		slog.Error("Error sending expiring products", "error", err.Error())
	}
}
