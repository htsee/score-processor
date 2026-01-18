package util_test

import (
	"github.com/htsee/score-processor/internal/util"
	"testing"
)

func TestCheckFiletype(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		filetype string
		want     string
	}{
		{"correct input", "../../testdata/test.pdf", "pdf", ""},
		{"Input not a PDF", "../../testdata/test.png", "pdf", "File \"../../testdata/test.png\" is not a pdf"},
		{"Input does not exist", "notExist.pdf", "pdf", "File \"notExist.pdf\" does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.CheckFileType(tt.input, tt.filetype)
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
