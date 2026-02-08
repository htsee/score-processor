package cmd

import (
	"fmt"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
)

var VSpliceCmd = &cobra.Command{
	Use:   "vsplice [inputs] [destination]",
	Short: "Splice 2 vertical images",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		if err := util.CheckValidIO(inputs, "png", destination); err != nil {
			return err
		}
		if err := VSplice(inputs, destination); err != nil {
			return err
		}
		return nil
	},
}

func VSplice(inputs []string, destination string) error {
	maxHeight := 0
	var pages []gocv.Mat
	index := 1
	var g errgroup.Group
	for i, input := range inputs {
		page := gocv.IMRead(input, gocv.IMReadGrayScale)
		if page.Empty() {
			return fmt.Errorf("cannot read image %q", input)
		}
		imgHeight := page.Rows()
		if imgHeight > maxHeight {
			maxHeight = imgHeight
		}

		pages = append(pages, page)
		if len(pages) >= 2 || i == len(inputs)-1 {
			pagesCopy := pages
			g.Go(func() error {
				return util.Combine(pagesCopy, maxHeight, index, "vertical", destination)
			})
			index++
			maxHeight = 0
			pages = pages[:0]
		}
	}
	return g.Wait()
}
