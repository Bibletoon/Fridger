package bot

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"image"
	"net/http"
)

func handlePhotoMatchFunc(upd *models.Update) bool {
	return len(upd.Message.Photo) > 0
}

func (b *BotWrapper) handlePhoto(ctx context.Context, upd *models.Update) {
	params := bot.GetFileParams{
		FileID: upd.Message.Photo[3].FileID,
	}

	file, err := b.bot.GetFile(ctx, &params)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	link := b.bot.FileDownloadLink(file)
	f, err := http.Get(link)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	defer f.Body.Close()

	img, _, _ := image.Decode(f.Body)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	reader := datamatrix.NewDataMatrixReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		fmt.Printf("Decode error\n")
		fmt.Printf("%v\n", err.Error())
		return
	}

	info, err := b.crptClient.GetByDatamatrix(ctx, result.GetText())
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   info.ProductName,
	}
	_, err = b.bot.SendMessage(ctx, &msg)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
}
