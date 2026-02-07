package cmd

import (
	"fmt"
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
		if err := util.CheckNonNegative(width, height); err != nil {
			return err
		}
		ratio := float64(width) / float64(height)
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Fit(input, destination, ratio)
		})
	},
}

func Fit(input, destination string, ratio float64) error {
	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	fitted, err := util.Fit(img, ratio)
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
