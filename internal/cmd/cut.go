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

var CutCmd = &cobra.Command{
	Use:   "cut [inputs] [destination]",
	Short: "Cut images into staves",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		for _, input := range inputs {
			if err := Cut(input, destination); err != nil {
				return err
			}
		}
		return nil
	},
}

func Cut(input string, destination string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	defer img.Close()

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	img_name, _, _ := strings.Cut(path.Base(input), ".")
	output_path := fmt.Sprintf("%s/%s_001.png", destination, img_name)

	gocv.IMWrite(output_path, img)
	return nil

}
