package cmd

import (
	"fmt"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
)

var SpliceCmd = &cobra.Command{
	Use:   "splice [inputs] [destination]",
	Short: "Splice images horizontally",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		if err := Splice(inputs, destination); err != nil {
			return err
		}
		return nil
	},
}

func Splice(inputs []string, destination string) error {
	maxWidth := 0
	currentHeight := 0
	var staves []gocv.Mat
	type staffGroup struct {
		end   int
		width int
	}
	var groups []staffGroup
	for i, input := range inputs {
		staff := gocv.IMRead(input, gocv.IMReadGrayScale)
		if staff.Empty() {
			return fmt.Errorf("cannot read image %q", input)
		}
		staves = append(staves, staff)
		imgWidth, imgHeight := staff.Cols(), staff.Rows()
		if imgWidth > maxWidth {
			maxWidth = imgWidth
		}
		currentHeight += imgHeight
		if i > 0 && float64(currentHeight) > float64(maxWidth)/(16.0/9.0) {
			groups = append(groups, staffGroup{end: i - 1, width: maxWidth})
			maxWidth = imgWidth
			currentHeight = imgHeight
		}
		if i == len(inputs)-1 {
			groups = append(groups, staffGroup{end: i, width: maxWidth})
		}
	}
	start := 0
	var g errgroup.Group
	for i := range groups {
		groupedStaves := staves[start : groups[i].end+1]
		g.Go(func() error {
			return util.Combine(groupedStaves, groups[i].width, i, "horizontal", destination)
		})
		start = groups[i].end + 1
	}

	return g.Wait()
}
