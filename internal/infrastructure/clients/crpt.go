package clients

import (
	"Fridger/internal/domain/interfaces/clients"
	"Fridger/internal/errors"
	"Fridger/internal/infrastructure/clients/dto"
	"context"
	"fmt"
	"github.com/perimeterx/marshmallow"
	"io"
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

func (c *crptClient) GetByDatamatrix(ctx context.Context, datamatrix string) (*dto.ProductInfoDto, error) {
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

	var info = dto.ProductInfoDto{}
	fields, err := marshmallow.Unmarshal(responseBytes, &info)
	if err != nil {
		return nil, err
	}

	if info.CodeFounded == false {
		return nil, errors.ErrProductNotFound
	}

	productData, ok := fields[info.Category+"Data"].(map[string]interface{})

	if !ok {
		return nil, errors.ErrResponseRead
	}

	expiration, ok := productData["expireDate"].(float64)
	if !ok {
		return nil, errors.ErrResponseRead
	}

	info.ExpirationDate = expiration

	return &info, nil
}

func addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Platform: iOS 17.2; AppVersion: 4.47.0; AppVersionCode: 7630; Device: iPhone 14 Pro;")
	req.Header.Add("Client", "iOS 17.2; AppVersion: 4.47.0; Device: iPhone 14 Pro;")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}
