package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var TrimCmd = &cobra.Command{
	Use:   "trim [inputs] [destination]",
	Short: "Trim image borders",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		top, err := cmd.Flags().GetInt("top")
		if err != nil {
			return err
		}
		bottom, err := cmd.Flags().GetInt("bottom")
		if err != nil {
			return err
		}
		left, err := cmd.Flags().GetInt("left")
		if err != nil {
			return err
		}
		right, err := cmd.Flags().GetInt("right")
		if err != nil {
			return err
		}
		if err := util.CheckNonNegative(top, bottom, left, right); err != nil {
			return err
		}
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Trim(input, destination, top, bottom, left, right)
		})
	},
}

func Trim(input, destination string, top, bottom, left, right int) error {
	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	trimmed := util.Trim(img, top, bottom, left, right)
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, trimmed)
	if err := trimmed.Close(); err != nil {
		return err
	}

	return nil
}
