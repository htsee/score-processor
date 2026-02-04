package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/htsee/score-processor/internal/cmd"
)

var pdfDir = "../../testdata/pdf"

var pageDir = "../../testdata/pages"

var cutDir = "../../testdata/cut"

func fileList(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		fileList = append(fileList, path)
	}
	return fileList, nil
}

var pdfList, _ = fileList(pdfDir)
var pageList, _ = fileList(pageDir)
var cutList, _ = fileList(cutDir)

func BenchmarkConvert(b *testing.B) {
	if err := cmd.ConvertBatch(pdfList, b.TempDir(), "1-N"); err != nil {
		b.Error(err)
	}
}

func BenchmarkCut(b *testing.B) {
	if err := cmd.CutBatch(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkDenoise(b *testing.B) {
	if err := cmd.DenoiseBatch(pageList, b.TempDir(), 2); err != nil {
		b.Error(err)
	}
}

func BenchmarkDeskew(b *testing.B) {
	if err := cmd.DeskewBatch(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkFit(b *testing.B) {
	if err := cmd.FitBatch(pageList, b.TempDir(), 16.0/9.0); err != nil {
		b.Error(err)
	}
}

func BenchmarkPad(b *testing.B) {
	if err := cmd.PadBatch(pageList, b.TempDir(), 10, 10); err != nil {
		b.Error(err)
	}
}

func BenchmarkRotate(b *testing.B) {
	if err := cmd.RotateBatch(pageList, b.TempDir(), 90); err != nil {
		b.Error(err)
	}
}

func BenchmarkSplice(b *testing.B) {
	if err := cmd.Splice(cutList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkTrim(b *testing.B) {
	if err := cmd.TrimBatch(pageList, b.TempDir(), 5, 5, 0, 0); err != nil {
		b.Error(err)
	}
}

func BenchmarkVsplice(b *testing.B) {
	if err := cmd.Splice(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}
