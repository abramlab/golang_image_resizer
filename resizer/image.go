package resizer

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

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
	return fmt.Sprintf("%s.jpg", i.filename)
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
	return fmt.Sprintf("%s.png", i.filename)
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
	return fmt.Sprintf("%s.gif", i.filename)
}

type Img struct {
	image.Image
	filename string
}

func (i *Img) resize(width, height uint) {
	i.Image = resize.Resize(width, height, i.Image, resize.Lanczos3)
}

// TODO: for future use
/*func (i Img) dimensions() (width int, height int) {
	bounds := i.Bounds()
	return bounds.Max.X, bounds.Max.Y
}*/
