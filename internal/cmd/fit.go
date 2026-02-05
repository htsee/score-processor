package cmd

import (
	"fmt"
	"image/color"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var FitCmd = &cobra.Command{
	Use:   "fit [inputs] [destination]",
	Short: "Fit image to an aspect ratio",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		width, err := cmd.Flags().GetInt("width")
		if err != nil {
			return err
		}
		height, err := cmd.Flags().GetInt("height")
		if err != nil {
			return err
		}
		ratio := float64(width) / float64(height)
		return util.Batch(inputs, func(input string) error {
			return FitCmdExecute(input, destination, ratio)
		})
	},
}

func FitCmdExecute(input, destination string, ratio float64) error {
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

	fitted, err := Fit(img, ratio)
	if err != nil {
		return fmt.Errorf("failed to fit image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, fitted)
	if err := fitted.Close(); err != nil {
		return err
	}

	return nil
}

func Fit(img gocv.Mat, ratio float64) (gocv.Mat, error) {
	w := float64(img.Cols())
	h := float64(img.Rows())
	fitted := gocv.NewMat()
	if w/h < ratio {
		padding := int((h*ratio - w) / 2)
		err := gocv.CopyMakeBorder(img, &fitted, 0, 0, padding, padding, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
		if err != nil {
			return img, err
		}
	} else {
		padding := int((w/ratio - h) / 2)
		err := gocv.CopyMakeBorder(img, &fitted, padding, padding, 0, 0, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
		if err != nil {
			return img, err
		}
	}
	return fitted, nil
}
