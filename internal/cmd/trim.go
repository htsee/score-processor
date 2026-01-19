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

	trimmed := Trim(img)
	img.Close()

	gocv.IMWrite(input, trimmed)
	trimmed.Close()

	return nil
}

func Trim(img gocv.Mat) gocv.Mat {
	trimSize := util.MmToPixel(10, img.Cols())
	trimmedRect := image.Rect(0, trimSize, img.Cols(), img.Rows()-trimSize)

	return img.Region(trimmedRect)
}
