package util

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func Pad(img gocv.Mat, vpad, hpad int) (gocv.Mat, error) {
	padded := gocv.NewMat()

	boundingBox, err := GetBoundingBox(img)
	if err != nil {
		return img, err
	}

	cropped := img.Region(boundingBox)

	err = gocv.CopyMakeBorder(cropped, &padded, vpad, vpad, hpad, hpad, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return img, err
	}
	if err := cropped.Close(); err != nil {
		return img, err
	}

	return padded, nil
}

func GetBoundingBox(img gocv.Mat) (image.Rectangle, error) {
	thresh := gocv.NewMat()
	gocv.Threshold(img, &thresh, 225, 255, gocv.ThresholdBinaryInv)

	nonZero := gocv.NewMat()
	if err := gocv.FindNonZero(thresh, &nonZero); err != nil {
		return image.Rectangle{}, err
	}
	if err := thresh.Close(); err != nil {
		return image.Rectangle{}, err
	}

	pointVector := gocv.NewPointVectorFromMat(nonZero)
	if err := nonZero.Close(); err != nil {
		return image.Rectangle{}, err
	}

	boundingRect := gocv.BoundingRect(pointVector)
	pointVector.Close()

	return boundingRect, nil
}
