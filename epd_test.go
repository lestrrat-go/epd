package epd_test

import (
	"os"
	"testing"

	"github.com/lestrrat-go/epd"
)

func TestEPD(t *testing.T) {
	if _, err := os.Stat("/dev/spidev0.0"); err != nil {
		t.Skip("Tests must be run on devices with /dev/spidev0.0")
		return
	}

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
