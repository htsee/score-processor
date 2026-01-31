package cmd

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var RotateCmd = &cobra.Command{
	Use:   "rotate [inputs] [destination]",
	Short: "Rotate images clockwise",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		angle, err := cmd.Flags().GetFloat64("angle")
		if err != nil {
			return err

		}
		for _, input := range inputs {
			if err := RotateCmdExecute(input, destination, angle); err != nil {
				return err
			}
		}
		return nil
	},
}

func RotateCmdExecute(input, destination string, angle float64) error {
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

	rotated, err := Rotate(img, -angle)
	if err != nil {
		return fmt.Errorf("failed to rotate image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, rotated)
	if err := rotated.Close(); err != nil {
		return err
	}

	return nil
}

func Rotate(img gocv.Mat, angle float64) (gocv.Mat, error) {
	rotated := gocv.NewMat()

	imgW, imgH := img.Cols(), img.Rows()

	centre := image.Point{imgW / 2, imgH / 2}

	rotationMatrix := gocv.GetRotationMatrix2D(centre, angle, 1.0)

	absCos := math.Abs(rotationMatrix.GetDoubleAt(0, 0))
	absSin := math.Abs(rotationMatrix.GetDoubleAt(0, 1))
	newW := int(float64(imgW)*absCos + float64(imgH)*absSin)
	newH := int(float64(imgH)*absCos + float64(imgW)*absSin)

	rotationMatrix.SetDoubleAt(0, 2, rotationMatrix.GetDoubleAt(0, 2)+float64(newW-imgW)/2.0)
	rotationMatrix.SetDoubleAt(1, 2, rotationMatrix.GetDoubleAt(1, 2)+float64(newH-imgH)/2.0)

	err := gocv.WarpAffineWithParams(img, &rotated, rotationMatrix, image.Point{newW, newH}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	if err != nil {
		return img, err
	}
	if err := rotationMatrix.Close(); err != nil {
		return img, err
	}

	return rotated, nil
}
