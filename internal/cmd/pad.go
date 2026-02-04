package cmd

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var PadCmd = &cobra.Command{
	Use:   "pad [inputs] [destination]",
	Short: "Pad image with white border",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		vpad, err := cmd.Flags().GetInt("vpad")
		if err != nil {
			return err
		}
		hpad, err := cmd.Flags().GetInt("hpad")
		if err != nil {
			return err
		}
		if err := PadBatch(inputs, destination, vpad, hpad); err != nil {
			return err
		}
		return nil
	},
}

func PadBatch(imgs []string, destination string, vpad, hpad int) error {
	for _, img := range imgs {
		if err := PadCmdExecute(img, destination, vpad, hpad); err != nil {
			return err
		}
	}
	return nil
}

func PadCmdExecute(input, destination string, vpad, hpad int) error {
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

	padded, err := Pad(img, vpad, hpad)
	if err != nil {
		return fmt.Errorf("failed to pad image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, padded)
	if err := padded.Close(); err != nil {
		return err
	}

	return nil
}

func Pad(img gocv.Mat, vpad, hpad int) (gocv.Mat, error) {
	padded := gocv.NewMat()

	boundingBox, err := getBoundingBox(img)
	if err != nil {
		return img, err
	}

	cropped := img.Region(boundingBox)

	err = gocv.CopyMakeBorder(cropped, &padded, vpad, vpad, hpad, hpad, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return img, err
	}
	if err := cropped.Close(); err != nil {
		return img, err
	}

	return padded, nil
}

func getBoundingBox(img gocv.Mat) (image.Rectangle, error) {
	thresh := gocv.NewMat()
	gocv.Threshold(img, &thresh, 225, 255, gocv.ThresholdBinaryInv)

	nonZero := gocv.NewMat()
	if err := gocv.FindNonZero(thresh, &nonZero); err != nil {
		return image.Rectangle{}, err
	}
	if err := thresh.Close(); err != nil {
		return image.Rectangle{}, err
	}

	pointVector := gocv.NewPointVectorFromMat(nonZero)
	if err := nonZero.Close(); err != nil {
		return image.Rectangle{}, err
	}

	boundingRect := gocv.BoundingRect(pointVector)
	pointVector.Close()

	return boundingRect, nil
}
