package clients

import (
	"Fridger/internal/domain/interfaces/clients"
	"Fridger/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type crptClient struct {
	httpClient http.Client
}

func NewCrptClient() clients.CrptClient {
	return &crptClient{
		httpClient: http.Client{},
	}
}

func (c *crptClient) GetByDatamatrix(ctx *context.Context, datamatrix string) (*models.ProductInfo, error) {
	link := fmt.Sprintf("https://mobile.api.crpt.ru/mobile/check?code=%s&codeType=datamatrix", url.QueryEscape(datamatrix))
	req, err := http.NewRequestWithContext(*ctx, "GET", link, nil)
	if err != nil {
		return nil, err
	}

	addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info = models.ProductInfo{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Platform: iOS 17.2; AppVersion: 4.47.0; AppVersionCode: 7630; Device: iPhone 14 Pro;")
	req.Header.Add("Client", "iOS 17.2; AppVersion: 4.47.0; Device: iPhone 14 Pro;")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}
