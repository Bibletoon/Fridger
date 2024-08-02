package handlers

import (
	"Fridger/internal/domain/interfaces/clients"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"image"
	"net/http"
)

type PhotoHandler struct {
	crptClient clients.CrptClient
}

func NewPhotoHandler(crptClient clients.CrptClient) *PhotoHandler {
	return &PhotoHandler{
		crptClient: crptClient,
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
	code, err := h.decodeImage(img)
	if err != nil {
		img = imaging.AdjustContrast(img, 50)
		code, err = h.decodeImage(img)
	}

	if err != nil {
		fmt.Printf("Decode error\n")
		fmt.Printf("%v\n", err.Error())
		return
	}

	info, err := h.crptClient.GetByDatamatrix(ctx, code)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   info.ProductName,
	}
	_, err = b.SendMessage(ctx, &msg)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
}

func (h *PhotoHandler) decodeImage(img image.Image) (string, error) {
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	reader := datamatrix.NewDataMatrixReader()
	result, err := reader.Decode(bmp, nil)

	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}
