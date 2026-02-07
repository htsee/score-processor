package cmd

import (
	"fmt"
	"os"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var SpliceCmd = &cobra.Command{
	Use:   "splice [inputs] [destination]",
	Short: "Splice images horizontally",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		if err := Splice(inputs, destination); err != nil {
			return err
		}
		return nil
	},
}

func Splice(inputs []string, destination string) error {
	for _, input := range inputs {
		if err := util.CheckFileType(input, "png"); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("cannot create folder %q: %w", destination, err)
	}

	maxWidth := 0
	currentHeight := 0
	var staves []gocv.Mat
	index := 1
	for i, input := range inputs {
		staff := gocv.IMRead(input, gocv.IMReadGrayScale)
		if staff.Empty() {
			return fmt.Errorf("cannot read image %q", input)
		}
		imgWidth, imgHeight := staff.Cols(), staff.Rows()
		if imgWidth > maxWidth {
			maxWidth = imgWidth
		}
		currentHeight += imgHeight
		if len(staves) != 0 && (float64(currentHeight) > float64(maxWidth)/(16.0/9.0)) {
			if err := util.Combine(staves, maxWidth, index, "horizontal", destination); err != nil {
				return err
			}
			index++
			maxWidth = imgWidth
			currentHeight = imgHeight
			staves = staves[:0]
		}

		staves = append(staves, staff)
		if i == len(inputs)-1 {
			if err := util.Combine(staves, maxWidth, index, "horizontal", destination); err != nil {
				return err
			}
		}
	}
	return nil
}
