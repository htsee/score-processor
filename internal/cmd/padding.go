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
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	padded := Padding(img)
	defer padded.Close()

	gocv.IMWrite(input, padded)

	return nil
}

func Padding(img gocv.Mat) gocv.Mat {
	padded := gocv.NewMat()
	paddingSize := img.Cols() / 50

	cropped := img.Region(getBoundingBox(img))
	defer cropped.Close()
	gocv.CopyMakeBorder(cropped, &padded, paddingSize, paddingSize, paddingSize, paddingSize, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	return padded
}

func getBoundingBox(img gocv.Mat) image.Rectangle {
	thresh := gocv.NewMat()
	defer thresh.Close()
	gocv.Threshold(img, &thresh, 50, 255, gocv.ThresholdBinaryInv)

	nonZero := gocv.NewMat()
	defer nonZero.Close()
	gocv.FindNonZero(thresh, &nonZero)

	pointVector := gocv.NewPointVectorFromMat(nonZero)
	defer pointVector.Close()

	return gocv.BoundingRect(pointVector)
}
