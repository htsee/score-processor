package util_test

import (
	"testing"

	"github.com/htsee/score-processor/internal/util"
	"gocv.io/x/gocv"
)

func TestCombine(t *testing.T) {
	testMat := gocv.IMRead("../../testdata/pages/test_001.png", gocv.IMReadGrayScale)
	testSlice := make([]gocv.Mat, 1)
	emptySlice := make([]gocv.Mat, 0)
	testSlice[0] = testMat
	tests := []struct {
		name        string
		input       []gocv.Mat
		orientation string
		want        string
	}{
		{"Correct input", testSlice, "vertical", ""},
		{"Invalid orientation", testSlice, "v", "invalid orientation"},
		{"Empty input", emptySlice, "vertical", "empty input"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.Combine(tt.input, 1000, 1, tt.orientation, t.TempDir())
			if err == nil {
				if tt.want != "" {
					t.Errorf("got nil, want %v", tt.want)
				}
				return
			}
			if err.Error() != tt.want {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}
