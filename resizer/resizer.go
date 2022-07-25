package resizer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Resizer struct {
	outPath       string
	width, height uint
}

func NewResizer(outPath string, opts ...Option) (*Resizer, error) {
	absPath, err := filepath.Abs(outPath)
	if err != nil {
		return nil, fmt.Errorf("creating absolute path: %w", err)
	}
	if err = os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory %q: %w", absPath, err)
	}
	r := &Resizer{outPath: absPath}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

func (r *Resizer) RunResizeWorker(in <-chan Image) {
	for img := range in {
		img.Resize(r.width, r.height)
		if err := saveImage(filepath.Join(r.outPath, img.Filename()), img); err != nil {
			fmt.Printf("saving resized image failed: %v\n", err)
			continue
		}
	}
}

func (r *Resizer) ResizeImageFile(path string) error {
	img, err := DecodeImageFile(path)
	if err != nil {
		return err
	}
	img.Resize(r.width, r.height)
	return saveImage(filepath.Join(r.outPath, img.Filename()), img)
}

func (r *Resizer) ResizeImage(reader io.Reader, filename string) error {
	img, err := DecodeImage(reader, filename)
	if err != nil {
		return err
	}
	img.Resize(r.width, r.height)
	return saveImage(filepath.Join(r.outPath, img.Filename()), img)
}

func saveImage(path string, img Image) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create image file %q: %w", path, err)
	}
	defer out.Close()

	return img.Encode(out)
}
