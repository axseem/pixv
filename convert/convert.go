package convert

import (
	"image"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

func RGBAtoColor(r, g, b, a uint32) uint32 {
	return (r/257)<<24 + (g/257)<<16 + (b/257)<<8 + (a / 257)
}

func OpenImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
