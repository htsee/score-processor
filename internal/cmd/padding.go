package cmd

import (
	"fmt"
	"image"
	"image/color"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var PaddingCmd = &cobra.Command{
	Use:   "padding [inputs]",
	Short: "Pad image with white border",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := paddingCmdExecute(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func paddingCmdExecute(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	padded, err := Padding(img)
	if err != nil {
		return fmt.Errorf("Failed to pad image: %w", err)
	}
	img.Close()

	gocv.IMWrite(input, padded)
	padded.Close()

	return nil
}

func Padding(img gocv.Mat) (gocv.Mat, error) {
	padded := gocv.NewMat()
	paddingSize := util.MmToPixel(5, img.Cols())

	boundingBox, err := getBoundingBox(img)
	if err != nil {
		return img, err
	}

	cropped := img.Region(boundingBox)

	err = gocv.CopyMakeBorder(cropped, &padded, paddingSize, paddingSize, paddingSize, paddingSize, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return img, err
	}
	cropped.Close()

	return padded, nil
}

func getBoundingBox(img gocv.Mat) (image.Rectangle, error) {
	thresh := gocv.NewMat()
	gocv.Threshold(img, &thresh, 50, 255, gocv.ThresholdBinaryInv)

	nonZero := gocv.NewMat()
	if err := gocv.FindNonZero(thresh, &nonZero); err != nil {
		return image.Rectangle{}, err
	}
	thresh.Close()

	pointVector := gocv.NewPointVectorFromMat(nonZero)
	nonZero.Close()

	boundingRect := gocv.BoundingRect(pointVector)
	pointVector.Close()

	return boundingRect, nil
}
