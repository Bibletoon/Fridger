package dto

type ProductInfoDto struct {
	Id              int64           `json:"id"`
	CodeFounded     bool            `json:"codeFounded"`
	CheckResult     bool            `json:"checkResult"`
	Category        string          `json:"category"`
	ProductName     string          `json:"productName"`
	CodeResolveData CodeResolveData `json:"codeResolveData"`
}

type CodeResolveData struct {
	Gtin string `json:"gtin"`
	Cis  string `json:"cis"`
}
