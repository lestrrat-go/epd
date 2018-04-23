package main

import (
	"image"
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
	e.SetFrameMemory(image.Black, 0, 0)
	e.DisplayFrame()
	e.SetFrameMemory(image.White, 0, 0)
	e.DisplayFrame()

	return nil
}
