package util

import (
	"image"

	"gocv.io/x/gocv"
)

func Trim(img gocv.Mat, top, bottom, left, right int) gocv.Mat {
	trimmedRect := image.Rect(left, top, img.Cols()-right, img.Rows()-bottom)

	return img.Region(trimmedRect)
}
