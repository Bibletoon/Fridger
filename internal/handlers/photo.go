package handlers

import (
	"Fridger/internal/domain/interfaces/services"
	"context"
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
	photosCount := len(upd.Message.Photo)
	params := bot.GetFileParams{
		FileID: upd.Message.Photo[photosCount-1].FileID,
	}

	file, err := b.GetFile(ctx, &params)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	link := b.FileDownloadLink(file)
	f, err := http.Get(link)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	defer f.Body.Close()

	img, _, _ := image.Decode(f.Body)
	code, err := h.photoService.DecodeDatamatrix(img)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	product, err := h.productService.AddProductByDatamatix(ctx, code)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	text := fmt.Sprintf("Product %s successfully added", product.Name)

	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   text,
	}
	_, err = b.SendMessage(ctx, &msg)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
}
