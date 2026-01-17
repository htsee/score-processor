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

	padded, err := Padding(img)
	if err != nil {
		return fmt.Errorf("Cannot pad image: %w", err)
	}
	defer padded.Close()

	gocv.IMWrite(input, padded)

	return nil
}

func Padding(img gocv.Mat) (gocv.Mat, error) {
	padded := gocv.NewMat()
	paddingSize := img.Cols() / 50

	boundingBox, err := getBoundingBox(img)
	if err != nil {
		return img, err
	}

	cropped := img.Region(boundingBox)
	defer cropped.Close()

	err = gocv.CopyMakeBorder(cropped, &padded, paddingSize, paddingSize, paddingSize, paddingSize, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return img, err
	}
	return padded, nil
}

func getBoundingBox(img gocv.Mat) (image.Rectangle, error) {
	thresh := gocv.NewMat()
	defer thresh.Close()
	gocv.Threshold(img, &thresh, 50, 255, gocv.ThresholdBinaryInv)

	nonZero := gocv.NewMat()
	defer nonZero.Close()
	if err := gocv.FindNonZero(thresh, &nonZero); err != nil {
		return image.Rectangle{}, err
	}

	pointVector := gocv.NewPointVectorFromMat(nonZero)
	defer pointVector.Close()

	return gocv.BoundingRect(pointVector), nil
}
