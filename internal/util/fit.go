package util

import (
	"image/color"

	"gocv.io/x/gocv"
)

func Fit(img gocv.Mat, ratio float64) (gocv.Mat, error) {
	w := float64(img.Cols())
	h := float64(img.Rows())
	fitted := gocv.NewMat()
	if w/h < ratio {
		padding := int((h*ratio - w) / 2)
		err := gocv.CopyMakeBorder(img, &fitted, 0, 0, padding, padding, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
		if err != nil {
			return img, err
		}
	} else {
		padding := int((w/ratio - h) / 2)
		err := gocv.CopyMakeBorder(img, &fitted, padding, padding, 0, 0, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
		if err != nil {
			return img, err
		}
	}
	return fitted, nil
}
