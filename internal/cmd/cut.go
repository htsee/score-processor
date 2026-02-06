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
		return util.Batch(inputs, func(input string) error {
			return Cut(input, destination)
		})
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
		return fmt.Errorf("cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	thresh := gocv.NewMat()

	gocv.Threshold(img, &thresh, 225, 255, gocv.ThresholdBinaryInv)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{3, 3})

	closed := gocv.NewMat()

	if err := gocv.MorphologyEx(thresh, &closed, gocv.MorphClose, kernel); err != nil {
		return err
	}

	if err := kernel.Close(); err != nil {
		return err
	}

	if err := thresh.Close(); err != nil {
		return err
	}

	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()
	numLabels := gocv.ConnectedComponentsWithStats(closed, &labels, &stats, &centroids)

	if err := closed.Close(); err != nil {
		return err
	}
	if err := centroids.Close(); err != nil {
		return err
	}

	minSizeForStaff := img.Cols() * 5

	var staves []staff

	for i := 1; i < numLabels; i++ {
		area := stats.GetIntAt(i, 4)

		if area >= int32(minSizeForStaff) {
			mask := gocv.NewMat()
			if err := gocv.InRangeWithScalar(
				labels,
				gocv.NewScalar(float64(i), 0, 0, 0),
				gocv.NewScalar(float64(i), 0, 0, 0),
				&mask,
			); err != nil {
				return err
			}

			top := stats.GetIntAt(i, 1)
			bottom := top + stats.GetIntAt(i, 3)

			staves = append(staves, staff{
				top:    top,
				bottom: bottom,
				mask:   mask.Clone(),
			})

			if err := mask.Close(); err != nil {
				return err
			}

		}
	}

	for i := 1; i < numLabels; i++ {
		area := stats.GetIntAt(i, 4)

		if area < int32(minSizeForStaff) {
			mask := gocv.NewMat()
			if err := gocv.InRangeWithScalar(
				labels,
				gocv.NewScalar(float64(i), 0, 0, 0),
				gocv.NewScalar(float64(i), 0, 0, 0),
				&mask,
			); err != nil {
				return err
			}

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

			if err := mask.Close(); err != nil {
				return err
			}
		}
	}
	if err := labels.Close(); err != nil {
		return err
	}

	if err := stats.Close(); err != nil {
		return err
	}

	for i, staff := range staves {
		boundingBox, err := util.GetBoundingBox(staff.mask)
		if err != nil {
			return err
		}

		maskCropped := staff.mask.Region(boundingBox)
		if err := staff.mask.Close(); err != nil {
			return err
		}

		imgCropped := img.Region(boundingBox)

		inverted := gocv.NewMat()
		if err := gocv.BitwiseNot(imgCropped, &inverted); err != nil {
			return err
		}
		if err := imgCropped.Close(); err != nil {
			return err
		}

		invertedMasked := gocv.NewMat()
		if err := inverted.CopyToWithMask(&invertedMasked, maskCropped); err != nil {
			return err
		}

		if err := maskCropped.Close(); err != nil {
			return err
		}

		masked := gocv.NewMat()
		if err := gocv.BitwiseNot(invertedMasked, &masked); err != nil {
			return err
		}
		if err := invertedMasked.Close(); err != nil {
			return err
		}

		padded, err := util.Pad(masked, 10, 10)
		if err != nil {
			return err
		}
		if err := masked.Close(); err != nil {
			return err
		}

		img_name, _ := strings.CutSuffix(path.Base(input), ".png")
		output_path := fmt.Sprintf("%s/%s_%03d.png", destination, img_name, i+1)

		gocv.IMWrite(output_path, padded)
		if err := padded.Close(); err != nil {
			return err
		}
	}
	if err := img.Close(); err != nil {
		return err
	}

	return nil
}
