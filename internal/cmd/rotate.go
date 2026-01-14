package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"os"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var RotateCmd = &cobra.Command{
	Use:   "rotate [inputs] [angle]",
	Short: "Rotate the image clockwise",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		angle := args[len(args)-1]
		for _, input := range inputs {
			if err := RotateCmdExecute(input, angle); err != nil {
				return err
			}
		}
		return nil
	},
}

func RotateCmdExecute(input, angle string) error {
	if path.Ext(input) != ".png" && path.Ext(input) != ".jpg" {
		return fmt.Errorf("File %q is not an image", input)
	}

	if _, err := os.Stat(input); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("File %q does not exist", input)
		}
		return fmt.Errorf("Cannot access file %q: %w", input, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	angleFloat, err := strconv.ParseFloat(angle, 64)

	if err != nil {
		return fmt.Errorf("Angle not valid: %w", err)
	}

	rotated := RotateImg(img, angleFloat)
	defer rotated.Close()

	gocv.IMWrite(input, rotated)

	return nil
}

func RotateImg(img gocv.Mat, angle float64) gocv.Mat {
	rotated := gocv.NewMat()

	imgW, imgH := img.Cols(), img.Rows()

	centre := image.Point{imgW / 2, imgH / 2}

	rotationMatrix := gocv.GetRotationMatrix2D(centre, angle, 1.0)
	defer rotationMatrix.Close()

	absCos := math.Abs(rotationMatrix.GetDoubleAt(0, 0))
	absSin := math.Abs(rotationMatrix.GetDoubleAt(0, 1))
	newW := int(float64(imgW)*absCos + float64(imgH)*absSin)
	newH := int(float64(imgH)*absCos + float64(imgW)*absSin)

	rotationMatrix.SetDoubleAt(0, 2, rotationMatrix.GetDoubleAt(0, 2)+float64(newW-imgW)/2.0)
	rotationMatrix.SetDoubleAt(1, 2, rotationMatrix.GetDoubleAt(1, 2)+float64(newH-imgH)/2.0)

	gocv.WarpAffineWithParams(img, &rotated, rotationMatrix, image.Point{newW, newH}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 0})

	return rotated
}
