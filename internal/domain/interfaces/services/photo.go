package services

import "image"

type PhotoService interface {
	DecodeDatamatrix(image image.Image) (string, error)
}
