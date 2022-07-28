package resizer

import (
	"fmt"
	"image"
	"io"
	"os"
	"strings"
)

var imageFormatToConstructor = map[string]func(i *Img) Image{
	"jpeg": func(i *Img) Image { return &JPEGImage{Img: i} },
	"png":  func(i *Img) Image { return &PNGImage{Img: i} },
	"gif":  func(i *Img) Image { return &GIFImage{Img: i} },
}

func DecodeImage(reader io.Reader, filename string) (Image, error) {
	decodedImage, format, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	constructor, ok := imageFormatToConstructor[format]
	if !ok {
		return nil, fmt.Errorf("unsupported image format %s", format)
	}
	return constructor(&Img{
		Image:    decodedImage,
		filename: filename,
	}), nil
}

func decodeImageFile(path string, rootPath string) (Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %s", path, err)
	}
	defer file.Close()

	return DecodeImage(file, strings.TrimPrefix(path, rootPath))
}
