package models

type ProductInfo struct {
	Id              int64           `json:"id"`
	CodeFounded     bool            `json:"codeFounded"`
	CheckResult     bool            `json:"checkResult"`
	Category        string          `json:"category"`
	ProductName     string          `json:"productName"`
	Expiration      int64           `json:"expiration"`
	CodeResolveData CodeResolveData `json:"codeResolveData"`
}

type CodeResolveData struct {
	Gtin string `json:"gtin"`
}
