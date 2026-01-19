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

var CutCmd = &cobra.Command{
	Use:   "cut [inputs] [destination]",
	Short: "Cut images into staves",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		for _, input := range inputs {
			if err := Cut(input, destination); err != nil {
				return err
			}
		}
		return nil
	},
}

type staff struct {
	top    int32
	bottom int32
	mask   gocv.Mat
}

func Cut(input string, destination string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	denoised, err := Denoise(img)
	if err != nil {
		return nil
	}
	img.Close()

	deskewed, err := Deskew(denoised)
	if err != nil {
		return err
	}
	denoised.Close()

	thresh := gocv.NewMat()

	gocv.Threshold(deskewed, &thresh, 220, 255, gocv.ThresholdBinaryInv)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{3, 3})

	closed := gocv.NewMat()

	if err := gocv.MorphologyEx(thresh, &closed, gocv.MorphClose, kernel); err != nil {
		return err
	}

	kernel.Close()
	thresh.Close()

	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()
	numLabels := gocv.ConnectedComponentsWithStats(closed, &labels, &stats, &centroids)
	closed.Close()
	centroids.Close()

	minSizeForStaff := deskewed.Cols() * 5

	var staves []staff

	for i := 1; i < numLabels; i++ {
		area := stats.GetIntAt(i, 4)

		if area >= int32(minSizeForStaff) {
			mask := gocv.NewMat()
			gocv.InRangeWithScalar(
				labels,
				gocv.NewScalar(float64(i), 0, 0, 0),
				gocv.NewScalar(float64(i), 0, 0, 0),
				&mask,
			)

			top := stats.GetIntAt(i, 1)
			bottom := top + stats.GetIntAt(i, 3)

			staves = append(staves, staff{
				top:    top,
				bottom: bottom,
				mask:   mask.Clone(),
			})

			mask.Close()
		}
	}

	for i := 1; i < numLabels; i++ {
		area := stats.GetIntAt(i, 4)

		if area < int32(minSizeForStaff) {
			mask := gocv.NewMat()
			gocv.InRangeWithScalar(
				labels,
				gocv.NewScalar(float64(i), 0, 0, 0),
				gocv.NewScalar(float64(i), 0, 0, 0),
				&mask,
			)

			top := stats.GetIntAt(i, 1)
			bottom := top + stats.GetIntAt(i, 3)

			closest := 0
			minDist := math.Inf(1)
			for i, staff := range staves {
				if top >= staff.top && bottom <= staff.bottom {
					closest = i
					break
				}
				dist := math.Min(math.Abs(float64(top-staff.top)), math.Abs(float64(bottom-staff.bottom)))
				if dist < minDist {
					closest = i
					minDist = dist
				}
			}

			if err := gocv.BitwiseOr(staves[closest].mask, mask, &staves[closest].mask); err != nil {
				return err
			}

			mask.Close()
		}
	}
	labels.Close()
	stats.Close()

	for i, staff := range staves {
		binary := gocv.NewMat()
		gocv.Threshold(staff.mask, &binary, 1, 255, gocv.ThresholdBinary)
		staff.mask.Close()

		nonZero := gocv.NewMat()
		if err := gocv.FindNonZero(binary, &nonZero); err != nil {
			return err
		}

		pointVector := gocv.NewPointVectorFromMat(nonZero)
		nonZero.Close()

		boundingRect := gocv.BoundingRect(pointVector)
		pointVector.Close()

		maskCropped := binary.Region(boundingRect)
		binary.Close()
		imgCropped := deskewed.Region(boundingRect)

		inverted := gocv.NewMat()
		if err := gocv.BitwiseNot(imgCropped, &inverted); err != nil {
			return err
		}
		imgCropped.Close()

		invertedMasked := gocv.NewMat()
		if err := inverted.CopyToWithMask(&invertedMasked, maskCropped); err != nil {
			return err
		}
		maskCropped.Close()

		masked := gocv.NewMat()
		if err := gocv.BitwiseNot(invertedMasked, &masked); err != nil {
			return err
		}
		invertedMasked.Close()

		padded, err := Padding(masked)
		if err != nil {
			return err
		}
		masked.Close()

		img_name, _ := strings.CutSuffix(path.Base(input), ".png")
		output_path := fmt.Sprintf("%s/%s_%03d.png", destination, img_name, i)

		gocv.IMWrite(output_path, padded)
		padded.Close()
	}
	deskewed.Close()
	return nil

}
