package cmd

import (
	"fmt"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var TrimCmd = &cobra.Command{
	Use:   "denoise [inputs]",
	Short: "Remove noise from image",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := trimCmdExecute(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func trimCmdExecute(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	trimmed, err := Trim(img)
	if err != nil {
		return fmt.Errorf("Failed to trim image: %w", err)
	}
	img.Close()

	gocv.IMWrite(input, trimmed)
	trimmed.Close()

	return nil
}

func Trim(img gocv.Mat) (gocv.Mat, error) {
	return img, nil
}
