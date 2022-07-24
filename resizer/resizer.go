package resizer

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type resizer struct {
	width, height uint
	in, out       string
	addPostfix    bool

	imagesCh chan *Image
	wg       *sync.WaitGroup

	resizedImages int
}

func newResizer(w, h uint, in, out string, addPostfix bool) *resizer {
	return &resizer{
		width:      w,
		height:     h,
		in:         in,
		out:        out,
		addPostfix: addPostfix,
		imagesCh:   make(chan *Image),
	}
}

func (r *resizer) scanDir() error {
	defer close(r.imagesCh)

	// TODO: recursive search by filepath.Walk
	files, err := ioutil.ReadDir(r.in)
	if err != nil {
		return fmt.Errorf("reading dir %q: %w", r.in, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		image, err := decodeImage(filepath.Join(r.in, file.Name()))
		if err != nil {
			fmt.Printf("skipping file: %v\n", err)
			continue
		}
		r.imagesCh <- image
	}
	return nil
}

func (r *resizer) resizeWorker() {
	defer r.wg.Done()
	for img := range r.imagesCh {
		img.resize(r.width, r.height)
		imageName := img.outFilename(r.addPostfix)
		if err := saveImage(filepath.Join(r.out, imageName), img); err != nil {
			fmt.Printf("saving resized image failed: %v\n", err)
			continue
		}
		r.resizedImages++
	}
}

func saveImage(path string, image *Image) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %q: %w", path, err)
	}
	defer out.Close()

	switch image.format {
	case "jpeg":
		return jpeg.Encode(out, image, nil)
	case "png":
		return png.Encode(out, image)
	case "gif":
		return gif.Encode(out, image, nil)
	default:
		return fmt.Errorf("unsupported image format: %s", image.format)
	}
}
