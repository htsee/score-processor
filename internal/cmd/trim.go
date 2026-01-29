package cmd

import (
	"fmt"
	"image"
	"os"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var TrimCmd = &cobra.Command{
	Use:   "trim [inputs] [destination]",
	Short: "Trim image borders",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
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
		for _, input := range inputs {
			if err := trimCmdExecute(input, destination, top, bottom, left, right); err != nil {
				return err
			}
		}
		return nil
	},
}

func trimCmdExecute(input, destination string, top, bottom, left, right int) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("cannot create folder %q: %w", destination, err)
	}

	img := gocv.IMRead(input, gocv.IMReadGrayScale)

	if img.Empty() {
		return fmt.Errorf("cannot read image %q", input)
	}

	trimmed := Trim(img, top, bottom, left, right)
	if err := img.Close(); err != nil {
		return err
	}

	img_name, _ := strings.CutSuffix(path.Base(input), ".png")
	output_path := fmt.Sprintf("%s/%s.png", destination, img_name)

	gocv.IMWrite(output_path, trimmed)
	if err := trimmed.Close(); err != nil {
		return err
	}

	return nil
}

func Trim(img gocv.Mat, top, bottom, left, right int) gocv.Mat {
	trimmedRect := image.Rect(left, top, img.Cols()-right, img.Rows()-bottom)

	return img.Region(trimmedRect)
}
