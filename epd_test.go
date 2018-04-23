package epd_test

import (
	"testing"

	"github.com/lestrrat-go/epd"
)

func TestEPD(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		e := epd.New()
		e.Reset()
	})
	t.Run("Clear", func(t *testing.T) {
		e := epd.New()
		e.Reset()
		e.ClearFrameMemory(0x00)
		e.DisplayFrame()
	})
}
