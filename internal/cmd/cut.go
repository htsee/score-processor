package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var CutCmd = &cobra.Command{
	Use:   "cut [inputs] [destination]",
	Short: "Cut the score into staves",
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
	if path.Ext(input) != ".png" && path.Ext(input) != ".jpg" {
		return fmt.Errorf("File %q is not an image", input)
	}

	if _, err := os.Stat(input); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("File %q does not exist", input)
		}
		return fmt.Errorf("Cannot access file %q: %w", input, err)
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	img_name, _, _ := strings.Cut(path.Base(input), ".")
	output_path := fmt.Sprintf("%s/%s_001.png", destination, img_name)

	gocv.IMWrite(output_path, img)
	return nil

}
