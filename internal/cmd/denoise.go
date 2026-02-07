package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var DenoiseCmd = &cobra.Command{
	Use:   "denoise [inputs] [destination]",
	Short: "Remove noise from image",
	Long:  "Remove noise from image. If elements are close to each other, they are considered part of a bigger element (so staccato dots and text would not be accidentally removed). Large size can be used to remove page numbers.",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		size, err := cmd.Flags().GetInt("size")
		if err != nil {
			return err
		}
		if err := util.CheckNonNegative(size); err != nil {
			return err
		}
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Denoise(input, destination, size)
		})
	},
}

func Denoise(input, destination string, size int) error {
	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	denoised, err := util.Denoise(img, size)
	if err != nil {
		return fmt.Errorf("failed to denoise image: %w", err)
	}
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, denoised)
	if err := denoised.Close(); err != nil {
		return err
	}

	return nil
}
