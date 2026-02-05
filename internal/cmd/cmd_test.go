package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/htsee/score-processor/internal/cmd"
	"github.com/htsee/score-processor/internal/util"
)

const pdfDir = "../../testdata/pdf"
const pageDir = "../../testdata/pages"
const cutDir = "../../testdata/cut"

var pdfList, pageList, cutList []string

func TestMain(m *testing.M) {
	var err error
	pdfList, err = util.FileList(pdfDir)
	if err != nil {
		fmt.Printf("failed to load pdfDir: %v", err)
		os.Exit(1)
	}
	pageList, err = util.FileList(pageDir)
	if err != nil {
		fmt.Printf("failed to load pageDir: %v", err)
		os.Exit(1)
	}
	cutList, err = util.FileList(cutDir)
	if err != nil {
		fmt.Printf("failed to load cutDir: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

const defaultPages = "1-N"

func BenchmarkConvertBatch(b *testing.B) {
	if err := util.Batch(pdfList, func(pdf string) error {
		return cmd.Convert(pdf, b.TempDir(), defaultPages)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkConvertSingle(b *testing.B) {
	if err := cmd.Convert(pdfList[0], b.TempDir(), defaultPages); err != nil {
		b.Error(err)
	}
}

func BenchmarkCutBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.Cut(page, b.TempDir())
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkCutSingle(b *testing.B) {
	if err := cmd.Cut(pageList[0], b.TempDir()); err != nil {
		b.Error(err)
	}
}

const defaultSize = 2

func BenchmarkDenoiseBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.DenoiseCmdExecute(page, b.TempDir(), defaultSize)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkDenoiseSingle(b *testing.B) {
	if err := cmd.DenoiseCmdExecute(pageList[0], b.TempDir(), defaultSize); err != nil {
		b.Error(err)
	}
}

func BenchmarkDeskewBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.DeskewCmdExecute(page, b.TempDir())
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkDeskewSingle(b *testing.B) {
	if err := cmd.DeskewCmdExecute(pageList[0], b.TempDir()); err != nil {
		b.Error(err)
	}
}

const defaultRatio = 16.0 / 9.0

func BenchmarkFitBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.FitCmdExecute(page, b.TempDir(), defaultRatio)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkFitSingle(b *testing.B) {
	if err := cmd.FitCmdExecute(pageList[0], b.TempDir(), defaultRatio); err != nil {
		b.Error(err)
	}
}

const defaultVpad, defaultHpad = 10, 10

func BenchmarkPadBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.PadCmdExecute(page, b.TempDir(), defaultVpad, defaultHpad)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkPadSingle(b *testing.B) {
	if err := cmd.PadCmdExecute(pageList[0], b.TempDir(), defaultVpad, defaultHpad); err != nil {
		b.Error(err)
	}
}

const defaultAngle = 90

func BenchmarkRotateBatch(b *testing.B) {
	if err := util.Batch(pageList, func(page string) error {
		return cmd.RotateCmdExecute(page, b.TempDir(), defaultAngle)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkRotateSingle(b *testing.B) {
	if err := cmd.RotateCmdExecute(pageList[0], b.TempDir(), defaultAngle); err != nil {
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
	if err := util.Batch(pageList, func(page string) error {
		return cmd.TrimCmdExecute(page, b.TempDir(), defaultTop, defaultBottom, defaultLeft, defaultRight)
	}); err != nil {
		b.Error(err)
	}
}

func BenchmarkTrimSingle(b *testing.B) {
	if err := cmd.TrimCmdExecute(pageList[0], b.TempDir(), defaultTop, defaultBottom, defaultLeft, defaultRight); err != nil {
		b.Error(err)
	}
}

func BenchmarkVsplice(b *testing.B) {
	if err := cmd.Splice(pageList, b.TempDir()); err != nil {
		b.Error(err)
	}
}
