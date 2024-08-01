package main

import (
	"Fridger/internal/infrastructure/clients"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"image"
	_ "image/jpeg"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New("")
	if err != nil {
		panic(err)
	}

	b.RegisterHandlerMatchFunc(func(upd *models.Update) bool {
		return len(upd.Message.Photo) > 0
	}, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		go handlePhoto(ctx, b, update)
	})

	b.Start(ctx)
}

func handlePhoto(ctx context.Context, b *bot.Bot, upd *models.Update) {
	params := bot.GetFileParams{
		FileID: upd.Message.Photo[3].FileID,
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

	crpt := clients.NewCrptClient()
	info, err := crpt.GetByDatamatrix(&ctx, result.GetText())
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
