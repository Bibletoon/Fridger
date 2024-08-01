package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"image"
	_ "image/jpeg"
	"net/http"
	"net/url"
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
	}, handlePhoto)

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

	info, err := getProductInfo(ctx, result.GetText())
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}

	msg := bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   info,
	}
	_, err = b.SendMessage(ctx, &msg)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
}

type productInfo struct {
	ProductName string `json:"productName"`
}

func getProductInfo(ctx context.Context, datamatrix string) (string, error) {
	client := &http.Client{}
	link := fmt.Sprintf("https://mobile.api.crpt.ru/mobile/check?code=%s&codeType=datamatrix", url.QueryEscape(datamatrix))
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Platform: iOS 17.2; AppVersion: 4.47.0; AppVersionCode: 7630; Device: iPhone 14 Pro;")
	req.Header.Add("Client", "iOS 17.2; AppVersion: 4.47.0; Device: iPhone 14 Pro;")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var info = productInfo{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", err
	}

	return info.ProductName, nil
}
