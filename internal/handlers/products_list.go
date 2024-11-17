package handlers

import (
	"Fridger/internal/domain/interfaces/handlers"
	"Fridger/internal/domain/interfaces/services"
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strings"
)

type productsListHandler struct {
	productService services.ProductService
}

func NewProductsListHandler(productService services.ProductService) handlers.Handler {
	return &productsListHandler{
		productService: productService,
	}
}

func (h *productsListHandler) Match(upd *models.Update) bool {
	return strings.HasPrefix(upd.Message.Text, "/list")
}

func (h *productsListHandler) Handle(ctx context.Context, b *bot.Bot, upd *models.Update) {
	message := h.handleInternal(ctx)
	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   message,
	}
	_, err := b.SendMessage(ctx, &msg)
	if err != nil {
		log.Printf("failed to send message: %e\n", err)
		return
	}
}

func (h *productsListHandler) handleInternal(ctx context.Context) string {
	products, err := h.productService.GetAllActiveProducts(ctx)

	if err != nil {
		log.Printf("failed to get all active products: %e\n", err)
		return fmt.Sprintf("Возникла ошибка при обработке команды: %e", err)
	}

	if len(products) == 0 {
		return "Продукты не найдены"
	}

	buf := bytes.Buffer{}
	buf.WriteString("Список продуктов:\n")

	for i, product := range products {
		buf.WriteString(fmt.Sprintf("%d) %s - %s\n", i, product.Name, product.ExpirationDate))
	}

	return buf.String()
}
