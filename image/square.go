package image

import (
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

func ResizeToSquare(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	size := img.Bounds().Size()
	if size.X == size.Y {
		// Cover is already square
		return path, nil
	}

	dim := size.Y
	if size.X > dim {
		dim = size.X
	}

	destImg := image.NewNRGBA(image.Rect(0, 0, dim, dim))

	draw.CatmullRom.Scale(destImg, destImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	tmpFile, err := os.CreateTemp("", "books-cover-*.jpg")
	if err != nil {
		return "", err
	}

	err = jpeg.Encode(tmpFile, destImg, nil)
	return tmpFile.Name(), err
}
