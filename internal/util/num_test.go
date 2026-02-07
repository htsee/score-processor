package util_test

import (
	"github.com/htsee/score-processor/internal/util"
	"testing"
)

func TestMmToPixel(t *testing.T) {
	got := util.MmToPixel(210, 1000)
	want := 1000
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCheckNonNegative(t *testing.T) {
	tests := []struct {
		name  string
		input int
		error bool
	}{
		{"Correct input", 0, false},
		{"Incorrect input", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.CheckNonNegative(tt.input)
			if err == nil && tt.error {
				t.Errorf("got nil, want error")
			}

			if err != nil && !tt.error {
				t.Errorf("got error %v, want nil", err)
			}
		})
	}
}
