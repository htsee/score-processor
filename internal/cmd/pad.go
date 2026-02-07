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

var PadCmd = &cobra.Command{
	Use:   "pad [inputs] [destination]",
	Short: "Pad image with white border",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		vpad, err := cmd.Flags().GetInt("vpad")
		if err != nil {
			return err
		}
		hpad, err := cmd.Flags().GetInt("hpad")
		if err != nil {
			return err
		}
		if err := util.CheckNonNegative(vpad, hpad); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Pad(input, destination, vpad, hpad)
		})
	},
}

func Pad(input, destination string, vpad, hpad int) error {
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

	padded, err := util.Pad(img, vpad, hpad)
	if err != nil {
		return fmt.Errorf("failed to pad image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, padded)
	if err := padded.Close(); err != nil {
		return err
	}

	return nil
}
