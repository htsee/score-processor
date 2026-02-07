package util_test

import (
	"os"
	"testing"

	"github.com/htsee/score-processor/internal/util"
)

var correct, incorrect, notExist []string

func TestMain(m *testing.M) {
	correct = append(correct, "../../testdata/pdf/test1.pdf")
	incorrect = append(incorrect, "../../testdata/pages/test_001.png")
	notExist = append(notExist, "notExist.pdf")
	os.Exit(m.Run())
}

func TestCheckFiletype(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		filetype string
		output   string
		want     string
	}{
		{"Correct input", correct, "pdf", "../../testdata", ""},
		{"Input not a PDF", incorrect, "pdf", "../../testdata", "file \"../../testdata/pages/test_001.png\" is not a pdf"},
		{"Input does not exist", notExist, "pdf", "../../testdata", "file \"notExist.pdf\" does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.CheckValidIO(tt.input, tt.filetype, tt.output)
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
