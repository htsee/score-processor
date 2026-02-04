package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/htsee/score-processor/internal/cmd"
)

const pdfDir = "../../testdata/pdf"

const pageDir = "../../testdata/pages"

const cutDir = "../../testdata/cut"

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

var pdfFirst = pdfList[0]
var pageFirst = pageList[0]

const defaultPages = "1-N"

func BenchmarkConvertBatch(b *testing.B) {
	if err := cmd.ConvertBatch(pdfList, b.TempDir(), defaultPages); err != nil {
		b.Error(err)
	}
}

func BenchmarkConvertSingle(b *testing.B) {
	if err := cmd.Convert(pdfFirst, b.TempDir(), defaultPages); err != nil {
		b.Error(err)
	}
}

func BenchmarkCutBatch(b *testing.B) {
	if err := cmd.CutBatch(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkCutSingle(b *testing.B) {
	if err := cmd.Cut(pageFirst, b.TempDir()); err != nil {
		b.Error(err)
	}
}

const defaultSize = 2

func BenchmarkDenoiseBatch(b *testing.B) {
	if err := cmd.DenoiseBatch(pageList, b.TempDir(), defaultSize); err != nil {
		b.Error(err)
	}
}

func BenchmarkDenoiseSingle(b *testing.B) {
	if err := cmd.DenoiseCmdExecute(pageFirst, b.TempDir(), defaultSize); err != nil {
		b.Error(err)
	}
}

func BenchmarkDeskewBatch(b *testing.B) {
	if err := cmd.DeskewBatch(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

func BenchmarkDeskewSingle(b *testing.B) {
	if err := cmd.DeskewCmdExecute(pageFirst, b.TempDir()); err != nil {
		b.Error(err)
	}
}

const defaultRatio = 16.0 / 9.0

func BenchmarkFitBatch(b *testing.B) {
	if err := cmd.FitBatch(pageList, b.TempDir(), defaultRatio); err != nil {
		b.Error(err)
	}
}

func BenchmarkFitSingle(b *testing.B) {
	if err := cmd.FitCmdExecute(pageFirst, b.TempDir(), defaultRatio); err != nil {
		b.Error(err)
	}
}

const defaultVpad, defaultHpad = 10, 10

func BenchmarkPadBatch(b *testing.B) {
	if err := cmd.PadBatch(pageList, b.TempDir(), defaultVpad, defaultHpad); err != nil {
		b.Error(err)
	}
}

func BenchmarkPadSingle(b *testing.B) {
	if err := cmd.PadCmdExecute(pageFirst, b.TempDir(), defaultVpad, defaultHpad); err != nil {
		b.Error(err)
	}
}

const defaultAngle = 90

func BenchmarkRotateBatch(b *testing.B) {
	if err := cmd.RotateBatch(pageList, b.TempDir(), defaultAngle); err != nil {
		b.Error(err)
	}
}

func BenchmarkRotateSingle(b *testing.B) {
	if err := cmd.RotateCmdExecute(pageFirst, b.TempDir(), defaultAngle); err != nil {
		b.Error(err)
	}
}

func BenchmarkSplice(b *testing.B) {
	if err := cmd.Splice(cutList, b.TempDir()); err != nil {
		b.Error(err)
	}
}

const defaultTop, defaultBottom, defaultLeft, defaultRight = 5, 5, 5, 5

func BenchmarkTrimBatch(b *testing.B) {
	if err := cmd.TrimBatch(pageList, b.TempDir(), defaultTop, defaultBottom, defaultLeft, defaultRight); err != nil {
		b.Error(err)
	}
}

func BenchmarkTrimSingle(b *testing.B) {
	if err := cmd.TrimCmdExecute(pageFirst, b.TempDir(), defaultTop, defaultBottom, defaultLeft, defaultRight); err != nil {
		b.Error(err)
	}
}

func BenchmarkVsplice(b *testing.B) {
	if err := cmd.Splice(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}
