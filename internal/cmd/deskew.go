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
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(img, &edges, 50, 200)

	lines := gocv.NewMat()
	defer lines.Close()
	gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/360, 100, float32(img.Cols())/4.0, 5)

	var angles []float64
	for i := 0; i < lines.Rows(); i++ {
		line := lines.GetVeciAt(i, 0)
		angle := math.Atan2(float64(line[3]-line[1]), float64(line[2]-line[0])) * (180.0 / math.Pi)
		angles = append(angles, angle)
	}

	medianAngle := angles[len(angles)/2]

	return Rotate(img, medianAngle)
}
