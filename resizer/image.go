package resizer

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"

	"github.com/nfnt/resize"
)

type Image interface {
	Encode(w io.Writer) error
	Resize(width, height uint)
	Filename() string
}

type JPEGImage struct {
	*Img
}

func (i *JPEGImage) Encode(w io.Writer) error {
	return jpeg.Encode(w, i.Image, nil)
}

func (i *JPEGImage) Resize(width, height uint) {
	i.resize(width, height)
}

func (i JPEGImage) Filename() string {
	if filepath.Ext(i.filename) == "" {
		return fmt.Sprintf("%s.jpg", i.filename)
	}
	return i.filename
}

type PNGImage struct {
	*Img
}

func (i *PNGImage) Encode(w io.Writer) error {
	return png.Encode(w, i.Image)
}

func (i *PNGImage) Resize(width, height uint) {
	i.resize(width, height)
}

func (i PNGImage) Filename() string {
	if filepath.Ext(i.filename) == "" {
		return fmt.Sprintf("%s.png", i.filename)
	}
	return i.filename
}

type GIFImage struct {
	*Img
}

func (i *GIFImage) Encode(w io.Writer) error {
	return gif.Encode(w, i.Image, nil)
}

func (i *GIFImage) Resize(width, height uint) {
	i.resize(width, height)
}

func (i GIFImage) Filename() string {
	if filepath.Ext(i.filename) == "" {
		return fmt.Sprintf("%s.gif", i.filename)
	}
	return i.filename
}

type Img struct {
	image.Image
	filename string
}

func (i *Img) resize(width, height uint) {
	i.Image = resize.Resize(width, height, i.Image, resize.Lanczos3)
}
