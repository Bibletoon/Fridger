package services

import (
	"Fridger/internal/domain/interfaces/services"
	"Fridger/internal/errors"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"image"
)

type photoService struct {
	dmReader *datamatrix.DataMatrixReader
}

func NewPhotoService(dmReader *datamatrix.DataMatrixReader) services.PhotoService {
	return &photoService{dmReader: dmReader}
}

func (s *photoService) DecodeDatamatrix(image image.Image) (string, error) {
	bmp, _ := gozxing.NewBinaryBitmapFromImage(image)
	reader := datamatrix.NewDataMatrixReader()
	result, err := reader.Decode(bmp, nil)

	if err != nil {
		return "", fmt.Errorf("%w: %w", errors.ErrCodeRead, err)
	}

	return result.GetText(), nil
}
