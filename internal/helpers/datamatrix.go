package helpers

import "strings"

func ParseCis(datamatrix string) string {
	idx := strings.Index(datamatrix[19:], "\u001D")
	return datamatrix[19 : 19+idx]
}
