package handlers

import (
	"Fridger/internal/domain/interfaces/services"
	models2 "Fridger/internal/domain/models"
	"Fridger/internal/errors"
	"context"
	errors2 "errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"image"
	"net/http"
)

type PhotoHandler struct {
	photoService   services.PhotoService
	productService services.ProductService
}

func NewPhotoHandler(photoService services.PhotoService, productService services.ProductService) *PhotoHandler {
	return &PhotoHandler{
		photoService:   photoService,
		productService: productService,
	}
}

func (h *PhotoHandler) Match(upd *models.Update) bool {
	return len(upd.Message.Photo) > 0
}

func (h *PhotoHandler) Handle(ctx context.Context, b *bot.Bot, upd *models.Update) {
	message := h.handleInternal(ctx, b, upd)
	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   message,
	}
	_, err := b.SendMessage(ctx, &msg)
	if err != nil {
		fmt.Printf("Error sending message %v\n", err.Error())
		return
	}
}

func (h *PhotoHandler) handleInternal(ctx context.Context, b *bot.Bot, upd *models.Update) string {
	img, err := extractImage(ctx, b, upd)
	if err != nil {
		return getMessage(nil, err)
	}
	code, err := h.photoService.DecodeDatamatrix(img)
	if err != nil {
		return getMessage(nil, err)
	}

	product, err := h.processCode(ctx, code)
	message := getMessage(product, err)
	return message
}

func extractImage(ctx context.Context, b *bot.Bot, upd *models.Update) (image.Image, error) {
	photosCount := len(upd.Message.Photo)
	params := bot.GetFileParams{
		FileID: upd.Message.Photo[photosCount-1].FileID,
	}

	file, err := b.GetFile(ctx, &params)
	if err != nil {
		return nil, err
	}

	link := b.FileDownloadLink(file)
	f, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer f.Body.Close()

	img, _, err := image.Decode(f.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (h *PhotoHandler) processCode(ctx context.Context, code string) (*models2.Product, error) {
	product, err := h.productService.FindProductByDatamatrix(ctx, code)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return h.productService.AddProductByDatamatix(ctx, code)
	}

	if product.IsActive {
		product.IsActive = false
		err = h.productService.DeleteProductByDatamatrix(ctx, code)
		return product, err
	}

	return product, errors.ErrProductExists
}

func getMessage(product *models2.Product, err error) string {
	if errors2.Is(err, errors.ErrProductExists) {
		return fmt.Sprintf("Product %s already exists", product.Name)
	}

	if err != nil {
		return fmt.Sprintf("Error processing product: %s", err.Error())
	}

	if product.IsActive {
		return fmt.Sprintf("Product %s successfully added", product.Name)
	}

	return fmt.Sprintf("Product %s successfully deleted", product.Name)
}
