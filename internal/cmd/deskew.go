package cmd

import (
	"fmt"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var DeskewCmd = &cobra.Command{
	Use:   "deskew [inputs]",
	Short: "Deskew images",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := deskewCmdExecute(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func deskewCmdExecute(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	deskewed := Deskew(img)
	defer deskewed.Close()

	gocv.IMWrite(input, deskewed)

	return nil
}

func Deskew(img gocv.Mat) gocv.Mat {
	return img
}
