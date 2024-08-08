package clients

import (
	"Fridger/internal/domain/interfaces/clients"
	"Fridger/internal/domain/models"
	errors2 "Fridger/internal/errors"
	"Fridger/internal/infrastructure/clients/dto"
	"context"
	"errors"
	"fmt"
	"github.com/perimeterx/marshmallow"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type crptClient struct {
	httpClient http.Client
}

func NewCrptClient() clients.CrptClient {
	return &crptClient{
		httpClient: http.Client{},
	}
}

func (c *crptClient) GetByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error) {
	link := fmt.Sprintf("https://mobile.api.crpt.ru/mobile/check?code=%s&codeType=datamatrix", url.QueryEscape(datamatrix))
	req, err := http.NewRequestWithContext(ctx, "GET", link, nil)
	if err != nil {
		return nil, err
	}

	addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productInfo = dto.ProductInfo{}
	fields, err := marshmallow.Unmarshal(responseBytes, &productInfo)
	if err != nil {
		return nil, err
	}

	if productInfo.CodeFounded == false {
		return nil, errors2.ErrProductNotFoundInCrpt
	}

	productData, ok := fields[productInfo.Category+"Data"].(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("field for category %s not found", productInfo.Category)
	}

	expiration, ok := productData["expireDate"].(float64)
	if !ok {
		return nil, errors.New("expireDate field not found")
	}

	gtin, err := strconv.ParseInt(productInfo.CodeResolveData.Gtin, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid gtin format %s", productInfo.CodeResolveData.Gtin)
	}

	product := models.Product{
		Name:           productInfo.ProductName,
		Gtin:           gtin,
		Serial:         productInfo.CodeResolveData.Serial,
		Category:       productInfo.Category,
		ExpirationDate: time.UnixMilli(int64(expiration)),
		IsActive:       true,
		CreatedAt:      time.Now(),
	}

	return &product, nil
}

func addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Platform: iOS 17.2; AppVersion: 4.47.0; AppVersionCode: 7630; Device: iPhone 14 Pro;")
	req.Header.Add("Client", "iOS 17.2; AppVersion: 4.47.0; Device: iPhone 14 Pro;")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}
