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
