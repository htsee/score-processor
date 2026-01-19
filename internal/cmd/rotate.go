package cmd

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var RotateCmd = &cobra.Command{
	Use:   "rotate [inputs] [angle]",
	Short: "Rotate images clockwise",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		angle := args[len(args)-1]
		for _, input := range inputs {
			if err := rotateCmdExecute(input, angle); err != nil {
				return err
			}
		}
		return nil
	},
}

func rotateCmdExecute(input, angle string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	angleFloat, err := strconv.ParseFloat(angle, 64)

	if err != nil {
		return fmt.Errorf("Angle not valid: %w", err)
	}

	rotated, err := Rotate(img, -angleFloat)
	if err != nil {
		return fmt.Errorf("Failed to rotate image: %w", err)
	}
	img.Close()

	gocv.IMWrite(input, rotated)
	rotated.Close()

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
	rotationMatrix.Close()

	return rotated, nil
}
