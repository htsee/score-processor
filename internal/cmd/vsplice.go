package cmd

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var VSpliceCmd = &cobra.Command{
	Use:   "vsplice [inputs] [destination]",
	Short: "Splice 2 vertical images",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		if err := VSplice(inputs, destination); err != nil {
			return err
		}
		return nil
	},
}

func VSplice(inputs []string, destination string) error {
	for _, input := range inputs {
		if err := util.CheckFileType(input, "png"); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	for i := 0; i < len(inputs); i += 2 {
		end := i + 2
		end = min(end, len(inputs))
		pair := inputs[i:end]
		if len(pair) >= 2 {
			img1 := gocv.IMRead(pair[0], gocv.IMReadGrayScale)
			img2 := gocv.IMRead(pair[1], gocv.IMReadGrayScale)
			if img1.Empty() {
				return fmt.Errorf("Cannot read image %q", pair[0])
			}
			if img2.Empty() {
				return fmt.Errorf("Cannot read image %q", pair[1])
			}

			spliced := gocv.NewMat()
			h1, h2 := img1.Rows(), img2.Rows()
			padding := math.Abs(float64(h1-h2)) / 2.0

			if h1 > h2 {
				gocv.CopyMakeBorder(img2, &img2, int(math.Ceil(padding)), int(math.Floor(padding)), 0, 0, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
			} else {
				gocv.CopyMakeBorder(img1, &img1, int(math.Ceil(padding)), int(math.Floor(padding)), 0, 0, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
			}

			if err := gocv.Hconcat(img1, img2, &spliced); err != nil {
				return err
			}
			img1.Close()
			img2.Close()

			fitted, err := Fit(spliced, 16.0/9.0)
			if err != nil {
				return err
			}
			spliced.Close()

			img_name, _ := strings.CutSuffix(path.Base(pair[0]), ".png")
			output_path := fmt.Sprintf("%s/%s.png", destination, img_name)
			gocv.IMWrite(output_path, fitted)
			fitted.Close()
		} else {
			img := gocv.IMRead(pair[0], gocv.IMReadGrayScale)
			if img.Empty() {
				return fmt.Errorf("Cannot read image %q", pair[0])
			}

			fitted, err := Fit(img, 16.0/9.0)
			if err != nil {
				return err
			}
			img.Close()

			img_name, _ := strings.CutSuffix(path.Base(pair[0]), ".png")
			output_path := fmt.Sprintf("%s/%s.png", destination, img_name)
			gocv.IMWrite(output_path, fitted)
			fitted.Close()
		}
	}
	return nil
}
