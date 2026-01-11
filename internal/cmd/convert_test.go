package cmd_test

import (
	"github.com/htsee/score-processor/internal/cmd"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		destination string
		want        string
	}{
		{"correct input", "../../testdata/test.pdf", t.TempDir(), ""},
		{"Input not a PDF", "../../testdata/test.png", t.TempDir(), "File \"../../testdata/test.png\" is not a PDF"},
		{"Destination not a folder", "../../testdata/test.pdf", "../../flake.lock", "\"../../flake.lock\" is not a folder"},
		{"Input does not exist", "notExist.pdf", t.TempDir(), "File \"notExist.pdf\" does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Convert(tt.input, tt.destination, "1-N")
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

func BenchmarkConvert(b *testing.B) {
	for b.Loop() {
		cmd.Convert("../../testdata/test.pdf", b.TempDir(), "1-N")
	}
}
