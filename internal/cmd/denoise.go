package cmd

import (
	"fmt"
	"image"
	"math"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var DenoiseCmd = &cobra.Command{
	Use:   "denoise [inputs] [destination]",
	Short: "Remove noise from image",
	Long:  "Remove noise from image. If elements are close to each other, they are considered part of a bigger element (so staccato dots and text would not be accidentally removed). Large size can be used to remove page numbers.",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		size, err := cmd.Flags().GetInt("size")
		if err != nil {
			return err
		}
		for _, input := range inputs {
			if err := DenoiseCmdExecute(input, destination, size); err != nil {
				return err
			}
		}
		return nil
	},
}

func DenoiseCmdExecute(input, destination string, size int) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	denoised, err := Denoise(img, size)
	if err != nil {
		return fmt.Errorf("failed to denoise image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, denoised)
	if err := denoised.Close(); err != nil {
		return err
	}

	return nil
}

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

	maxSizeForNoise := math.Pow(float64(util.MmToPixel(size, img.Cols())), 2)

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
