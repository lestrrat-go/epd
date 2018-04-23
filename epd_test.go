package epd_test

import (
	"testing"

	"github.com/lestrrat-go/epd"
)

func TestEPD(t *testing.T) {
	e := epd.New()
	t.Run("Reset", func(t *testing.T) {
		e.Reset()
		e.WaitUntilIdle(nil)
	})
	t.Run("Clear", func(t *testing.T) {
		e.ClearFrameMemory(0x00)
		e.DisplayFrame()
	})
}
