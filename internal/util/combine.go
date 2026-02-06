package util

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func Combine(imgs []gocv.Mat, length, index int, orientation, destination string) error {
	if orientation != "horizontal" && orientation != "vertical" {
		return errors.New("invalid orientation")
	}
	current := imgs[0]
	if orientation == "horizontal" {
		for i, img := range imgs {
			staffWidth := img.Cols()
			if staffWidth < length {
				padding := float64(length-staffWidth) / 2
				err := gocv.CopyMakeBorder(img, &img, 0, 0, int(math.Ceil(padding)), int(math.Floor(padding)), gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
				if err != nil {
					return err
				}
			}
			if i > 0 {
				if err := gocv.Vconcat(current, img, &current); err != nil {
					return err
				}
				if err := img.Close(); err != nil {
					return err
				}
			}
		}
	} else {
		for i, img := range imgs {
			staffHeight := img.Rows()
			if staffHeight < length {
				padding := float64(length-staffHeight) / 2
				err := gocv.CopyMakeBorder(img, &img, int(math.Ceil(padding)), int(math.Floor(padding)), 0, 0, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
				if err != nil {
					return err
				}
			}
			if i > 0 {
				if err := gocv.Hconcat(current, img, &current); err != nil {
					return err
				}
				if err := img.Close(); err != nil {
					return err
				}
			}
		}
	}
	fitted, err := Fit(current, 16.0/9.0)
	if err != nil {
		return err
	}
	if err := current.Close(); err != nil {
		return err
	}
	output_path := fmt.Sprintf("%s/%03d.png", destination, index)
	gocv.IMWrite(output_path, fitted)
	if err := fitted.Close(); err != nil {
		return err
	}
	return nil
}
