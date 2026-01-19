package cmd

import (
	"fmt"
	"math"

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

	if img.Empty() {
		return fmt.Errorf("Cannot read image %q", input)
	}

	deskewed, err := Deskew(img)
	if err != nil {
		return fmt.Errorf("Cannot deskew image: %w", err)
	}
	img.Close()

	gocv.IMWrite(input, deskewed)
	deskewed.Close()

	return nil
}

func Deskew(img gocv.Mat) (gocv.Mat, error) {
	edges := gocv.NewMat()
	if err := gocv.Canny(img, &edges, 50, 200); err != nil {
		return img, err
	}

	lines := gocv.NewMat()

	if err := gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/360, 100, float32(img.Cols())/4.0, 5); err != nil {
		return img, err
	}
	edges.Close()

	var angles []float64
	for i := 0; i < lines.Rows(); i++ {
		line := lines.GetVeciAt(i, 0)
		angle := math.Atan2(float64(line[3]-line[1]), float64(line[2]-line[0])) * (180.0 / math.Pi)
		angles = append(angles, angle)
	}
	lines.Close()

	medianAngle := angles[len(angles)/2]

	rotated, err := Rotate(img, medianAngle)
	if err != nil {
		return img, err
	}

	return rotated, nil
}
