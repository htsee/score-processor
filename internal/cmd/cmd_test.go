package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/htsee/score-processor/internal/cmd"
)

// var pdfDir = "../../testdata"

var pageDir = "../../testdata/pages"

var cutDir = "../../testdata/cut"

func BenchmarkConvert(b *testing.B) {
	if err := cmd.Convert("../../testdata/test.pdf", b.TempDir(), "1-N"); err != nil {
		b.Error(err)
	}
}

func BenchmarkCut(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.Cut(imgPath, b.TempDir()); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDenoise(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.DenoiseCmdExecute(imgPath, b.TempDir(), 2); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDeskew(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.DeskewCmdExecute(imgPath, b.TempDir()); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFit(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.FitCmdExecute(imgPath, b.TempDir(), 16.0/9.0); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkPad(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.PadCmdExecute(imgPath, b.TempDir(), 10, 10); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRotate(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.RotateCmdExecute(imgPath, b.TempDir(), 90); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSplice(b *testing.B) {
	imgs, err := os.ReadDir(cutDir)
	if err != nil {
		b.Error(err)
	}
	var imgsPath []string
	for _, img := range imgs {
		imgPath := filepath.Join(cutDir, img.Name())
		imgsPath = append(imgsPath, imgPath)
	}
	if err := cmd.Splice(imgsPath, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkTrim(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		if err := cmd.TrimCmdExecute(imgPath, b.TempDir(), 5, 5, 0, 0); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkVsplice(b *testing.B) {
	imgs, err := os.ReadDir(pageDir)
	if err != nil {
		b.Error(err)
	}
	var imgsPath []string
	for _, img := range imgs {
		imgPath := filepath.Join(pageDir, img.Name())
		imgsPath = append(imgsPath, imgPath)
	}
	if err := cmd.Splice(imgsPath, b.TempDir()); err != nil {
		b.Error(err)
	}
}
