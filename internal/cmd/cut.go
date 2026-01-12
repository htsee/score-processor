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
	Use:   "cut [input] [destination]",
	Short: "Cut the score into staves",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		destination := args[1]
		if err := Cut(input, destination); err != nil {
			return err
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

	pdf_name, _, _ := strings.Cut(path.Base(input), ".")
	output_path := fmt.Sprintf("%s/%s_001.png", destination, pdf_name)

	img := gocv.IMRead(input, gocv.IMReadGrayScale)
	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}
	gocv.IMWrite(output_path, img)
	return nil

}
