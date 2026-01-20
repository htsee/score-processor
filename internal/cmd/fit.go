package cmd

import (
	"fmt"
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
		for _, input := range inputs {
			if err := fitCmdExecute(input, destination, ratio); err != nil {
				return err
			}
		}
		return nil
	},
}

func fitCmdExecute(input, destination string, ratio float64) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	fitted, err := Fit(img, ratio)
	if err != nil {
		return fmt.Errorf("Failed to denoise image: %w", err)
	}
	img.Close()

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, fitted)
	fitted.Close()

	return nil
}

func Fit(img gocv.Mat, ratio float64) (gocv.Mat, error) {
	fitted := gocv.NewMat()
	return fitted, nil
}
