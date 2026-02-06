package util

import (
	"image"
	"math"

	"gocv.io/x/gocv"
)

func Denoise(img gocv.Mat, size int) (gocv.Mat, error) {
	denoised := gocv.NewMat()

	thresh := gocv.NewMat()

	gocv.Threshold(img, &thresh, 225, 255, gocv.ThresholdBinaryInv)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{11, 11})

	closed := gocv.NewMat()

	if err := gocv.MorphologyEx(thresh, &closed, gocv.MorphClose, kernel); err != nil {
		return img, err
	}

	if err := thresh.Close(); err != nil {
		return img, err
	}

	if err := kernel.Close(); err != nil {
		return img, err
	}

	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()
	numLabels := gocv.ConnectedComponentsWithStats(closed, &labels, &stats, &centroids)
	if err := closed.Close(); err != nil {
		return img, err
	}

	if err := centroids.Close(); err != nil {
		return img, err
	}

	maxSizeForNoise := math.Pow(float64(MmToPixel(size, img.Cols())), 2)

	mergedMask := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	mergedMask.SetTo(gocv.Scalar{Val1: 255})

	for i := 1; i < numLabels; i++ {
		area := stats.GetIntAt(i, 4)

		if area <= int32(maxSizeForNoise) {
			mask := gocv.NewMat()
			err := gocv.InRangeWithScalar(
				labels,
				gocv.NewScalar(float64(i), 0, 0, 0),
				gocv.NewScalar(float64(i), 0, 0, 0),
				&mask,
			)

			if err != nil {
				return img, err
			}

			if err := gocv.BitwiseXor(mergedMask, mask, &mergedMask); err != nil {
				return img, err
			}
			if err := mask.Close(); err != nil {
				return img, err
			}
		}
	}
	if err := labels.Close(); err != nil {
		return img, err
	}

	if err := stats.Close(); err != nil {
		return img, err
	}

	inverted := gocv.NewMat()

	if err := gocv.BitwiseNot(img, &inverted); err != nil {
		return img, err
	}

	invertedDenoised := gocv.NewMat()

	if err := inverted.CopyToWithMask(&invertedDenoised, mergedMask); err != nil {
		return img, err
	}
	if err := inverted.Close(); err != nil {
		return img, err
	}
	if err := mergedMask.Close(); err != nil {
		return img, err
	}
	if err := gocv.BitwiseNot(invertedDenoised, &denoised); err != nil {
		return img, err
	}
	if err := invertedDenoised.Close(); err != nil {
		return img, err
	}

	return denoised, nil
}
