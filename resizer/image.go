package resizer

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	image.Image
	name   string
	format string
}

func (i *Image) resize(width, height uint) {
	i.Image = resize.Resize(width, height, i.Image, resize.Lanczos3)
}

func (i Image) dimensions() (width int, height int) {
	bounds := i.Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func (i Image) outFilename(addPostfix bool) string {
	if addPostfix {
		w, h := i.dimensions()
		return fmt.Sprintf("%s_%dx%d%s", i.name, w, h, i.ext())
	}
	return fmt.Sprintf("%s%s", i.name, i.ext())
}

func (i Image) ext() string {
	var ext string
	switch i.format {
	case "jpeg":
		ext = ".jpg"
	case "png":
		ext = ".png"
	case "gif":
		ext = ".gif"
	}
	return ext
}

func decodeImage(path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %s", path, err)
	}
	defer file.Close()

	base, ext := filepath.Base(file.Name()), filepath.Ext(file.Name())
	decodedImage, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode file %s: %s", path, err)
	}
	return &Image{
		Image:  decodedImage,
		name:   strings.TrimSuffix(base, ext),
		format: format,
	}, nil
}
