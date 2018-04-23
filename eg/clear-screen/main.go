package main

import (
	"log"
	"os"

	"github.com/lestrrat-go/epd"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

func _main() error {
	e := epd.New()
	e.ClearFrameMemory(0x00)
	e.DisplayFrame()

	return nil
}
