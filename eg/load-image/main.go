package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/lestrrat-go/epd"
	"github.com/pkg/errors"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

func _main() error {
	if len(os.Args) < 2 {
		return errors.New(`usage: go run eg/load-image/main.go <file.bmp>`)
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return errors.Wrapf(err, `failed to open file %s`, os.Args[1])
	}

	rdr := bytes.NewReader(buf)

	_, format, err := image.DecodeConfig(rdr)
	if err != nil {
		return errors.Wrap(err, `failed to decode image config`)
	}

	var decoder func(io.Reader) (image.Image, error)
	switch format {
	case "png":
		decoder = png.Decode
	case "gif":
		decoder = gif.Decode
	case "jpeg":
		decoder = jpeg.Decode
	default:
		return errors.Errorf(`unsupported image type %s`, format)
	}

	rdr.Seek(0, 0)
	im, err := decoder(rdr)
	if err != nil {
		return errors.Wrap(err, `failed to decode image`)
	}

	e := epd.New()
	e.SetFrameMemory(im, 0, 0)
	e.DisplayFrame()

	return nil
}
