package util

import (
	"math"

	"gocv.io/x/gocv"
)

func Deskew(img gocv.Mat) (gocv.Mat, error) {
	edges := gocv.NewMat()
	if err := gocv.Canny(img, &edges, 50, 200); err != nil {
		return img, err
	}

	lines := gocv.NewMat()

	if err := gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/360, 100, float32(img.Cols())/4.0, 5); err != nil {
		return img, err
	}
	if err := edges.Close(); err != nil {
		return img, err
	}

	var angles []float64
	for i := 0; i < lines.Rows(); i++ {
		line := lines.GetVeciAt(i, 0)
		angle := math.Atan2(float64(line[3]-line[1]), float64(line[2]-line[0])) * (180.0 / math.Pi)
		angles = append(angles, angle)
	}
	if err := lines.Close(); err != nil {
		return img, err
	}

	medianAngle := angles[len(angles)/2]

	rotated, err := Rotate(img, medianAngle)
	if err != nil {
		return img, err
	}

	return rotated, nil
}
