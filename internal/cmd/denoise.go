package cmd

import (
	"fmt"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var DenoiseCmd = &cobra.Command{
	Use:   "denoise [inputs]",
	Short: "Remove noise from image",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := denoiseCmdExecute(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func denoiseCmdExecute(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	denoised := Denoise(img)
	defer denoised.Close()

	gocv.IMWrite(input, denoised)

	return nil
}

func Denoise(img gocv.Mat) gocv.Mat {
	denoised := gocv.NewMat()
	return denoised
}
