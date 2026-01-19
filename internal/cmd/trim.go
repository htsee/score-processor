package cmd

import (
	"fmt"
	"image"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var TrimCmd = &cobra.Command{
	Use:   "trim [inputs]",
	Short: "Trim image borders",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
		for _, input := range args {
			if err := trimCmdExecute(input, top, bottom, left, right); err != nil {
				return err
			}
		}
		return nil
	},
}

func trimCmdExecute(input string, top, bottom, left, right int) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	trimmed := Trim(img, top, bottom, left, right)
	img.Close()

	gocv.IMWrite(input, trimmed)
	trimmed.Close()

	return nil
}

func Trim(img gocv.Mat, top, bottom, left, right int) gocv.Mat {
	trimmedRect := image.Rect(left, top, img.Cols()-right, img.Rows()-bottom)

	return img.Region(trimmedRect)
}
