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
		{"Correct input", "../../testdata/pdf/test1.pdf", "pdf", ""},
		{"Input not a PDF", "test.png", "pdf", "file \"test.png\" is not a pdf"},
		{"Input does not exist", "notExist.pdf", "pdf", "file \"notExist.pdf\" does not exist"},
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

func TestFileList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		error bool
	}{
		{"Correct input", "../../testdata/pdf/", false},
		{"Input does not exist", "notExist", true},
		{"Input is not a folder", "../../testdata/pdf/test1.pdf", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := util.FileList(tt.input)
			if err == nil && tt.error {
				t.Errorf("got nil, want error")
			}

			if err != nil && !tt.error {
				t.Errorf("got error %v, want nil", err)
			}
		})
	}
}
