package cmd

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var SpliceCmd = &cobra.Command{
	Use:   "splice [inputs] [destination]",
	Short: "Splice images",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		if err := Splice(inputs, destination); err != nil {
			return err
		}
		return nil
	},
}

func Splice(inputs []string, destination string) error {
	for _, input := range inputs {
		if err := util.CheckFileType(input, "png"); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	maxWidth := 0
	currentHeight := 0
	var image []gocv.Mat
	index := 0
	for i, input := range inputs {
		img := gocv.IMRead(input, gocv.IMReadGrayScale)
		if img.Empty() {
			return fmt.Errorf("Cannot read image %q", input)
		}
		imgWidth, imgHeight := img.Cols(), img.Rows()
		if imgWidth > maxWidth {
			maxWidth = imgWidth
		}
		currentHeight += imgHeight
		if len(image) != 0 && (float64(currentHeight) > float64(maxWidth)/(16.0/9.0) || i == len(inputs)-1) {
			if err := combine(image, maxWidth, index, destination); err != nil {
				return err
			}
			index++
			maxWidth = imgWidth
			currentHeight = imgHeight
			image = image[:0]
		}

		image = append(image, img)
	}
	return nil
}

func combine(image []gocv.Mat, width, index int, destination string) error {
	current := image[0]
	for i, staff := range image {
		staffWidth := staff.Cols()
		if staffWidth < width {
			padding := float64(width-staffWidth) / 2
			gocv.CopyMakeBorder(staff, &staff, 0, 0, int(math.Ceil(padding)), int(math.Floor(padding)), gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
		}
		if i > 0 {
			if err := gocv.Vconcat(current, staff, &current); err != nil {
				return err
			}
			staff.Close()
		}
	}
	fitted, err := Fit(current, 16.0/9.0)
	if err != nil {
		return err
	}
	current.Close()
	output_path := fmt.Sprintf("%s/%03d.png", destination, index)
	gocv.IMWrite(output_path, fitted)
	fitted.Close()
	return nil
}
