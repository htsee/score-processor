package cmd

import (
	"fmt"
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
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Rotate(input, destination, angle)
		})
	},
}

func Rotate(input, destination string, angle float64) error {
	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	rotated, err := util.Rotate(img, -angle)
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
