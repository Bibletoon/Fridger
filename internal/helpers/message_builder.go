package helpers

import (
	"Fridger/internal/domain/models"
	"bytes"
	"fmt"
)

func BuildExpiredProductsListMessage(products *models.ProductsCollection) string {
	if products.Len() == 0 {
		return ""
	}

	buf := bytes.Buffer{}

	expiringProducts := products.GetNotExpired()
	if len(expiringProducts) > 0 {
		buf.WriteString("<b>Продукты, которые скоро испортятся:</b>\n\n")

		for i, product := range expiringProducts {
			buf.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, product.Name, product.ExpirationDate.Format("02 January 2006")))
		}
	}

	expiredProducts := products.GetExpired()

	if len(expiredProducts) > 0 {
		if len(expiringProducts) > 0 {
			buf.WriteString("\n")
		}

		buf.WriteString("<b>Продукты, которые уже испортились:</b>\n\n")

		for i, product := range expiredProducts {
			buf.WriteString(fmt.Sprintf("%d. %s - %s", i+1, product.Name, product.ExpirationDate.Format("02 January 2006")))
		}
	}

	return buf.String()
}
