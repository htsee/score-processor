package cmd

import (
	"fmt"
	"image"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var DenoiseCmd = &cobra.Command{
	Use:   "denoise [inputs]",
	Short: "Remove noise from image",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := denoiseCmdExecute(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func denoiseCmdExecute(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	denoised, err := Denoise(img)
	if err != nil {
		return fmt.Errorf("Failed to denoise image: %w", err)
	}
	defer denoised.Close()

	gocv.IMWrite(input, denoised)

	return nil
}

func Denoise(img gocv.Mat) (gocv.Mat, error) {
	denoised := gocv.NewMat()

	thresh := gocv.NewMat()
	defer thresh.Close()

	gocv.Threshold(img, &thresh, 225, 255, gocv.ThresholdBinaryInv)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{5, 5})
	defer kernel.Close()

	closed := gocv.NewMat()
	defer closed.Close()

	if err := gocv.MorphologyEx(thresh, &closed, gocv.MorphClose, kernel); err != nil {
		return img, err
	}

	labels := gocv.NewMat()
	defer labels.Close()
	stats := gocv.NewMat()
	defer stats.Close()
	centroids := gocv.NewMat()
	defer centroids.Close()
	numLabels := gocv.ConnectedComponentsWithStats(closed, &labels, &stats, &centroids)

	maxSizeForNoise := img.Cols() / 500

	mergedMask := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	mergedMask.SetTo(gocv.Scalar{Val1: 255})
	defer mergedMask.Close()

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
		}
	}

	inverted := gocv.NewMat()
	defer inverted.Close()

	if err := gocv.BitwiseNot(img, &inverted); err != nil {
		return img, err
	}

	invertedDenoised := gocv.NewMat()
	defer invertedDenoised.Close()

	if err := inverted.CopyToWithMask(&invertedDenoised, mergedMask); err != nil {
		return img, err
	}

	if err := gocv.BitwiseNot(invertedDenoised, &denoised); err != nil {
		return img, err
	}

	return denoised, nil
}
