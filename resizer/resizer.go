package resizer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

type Resizer struct {
	input, output string
	width, height uint
	workersNum    int

	resizedImages uint32
}

func NewResizer(input, output string, opts ...Option) (*Resizer, error) {
	if _, err := os.Stat(input); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%q does not exist", input)
	}
	absOutput, err := filepath.Abs(output)
	if err != nil {
		return nil, fmt.Errorf("output absolute path: %w", err)
	}
	if err = os.MkdirAll(absOutput, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory %q: %w", absOutput, err)
	}
	r := &Resizer{
		input:      input,
		output:     absOutput,
		workersNum: runtime.NumCPU(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

type RunStat struct {
	ResizedImages uint
}

func (r *Resizer) Run(ctx context.Context) (*RunStat, error) {
	imagesCh, err := scanDir(r.input)
	if err != nil {
		return nil, fmt.Errorf("scanning dir: %w", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < r.workersNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.RunResizeWorker(ctx, imagesCh)
		}()
	}

	wg.Wait()
	return &RunStat{ResizedImages: uint(r.resizedImages)}, nil
}

func (r *Resizer) RunResizeWorker(ctx context.Context, in <-chan Image) {
	for {
		select {
		case <-ctx.Done():
			return
		case img, ok := <-in:
			if !ok {
				return
			}
			img.Resize(r.width, r.height)
			if err := saveImage(filepath.Join(r.output, img.Filename()), img); err != nil {
				fmt.Printf("saving resized image failed: %v\n", err)
				continue
			}
			atomic.AddUint32(&r.resizedImages, 1)
		}
	}
}

func (r *Resizer) ResizeImageFile(path string) error {
	img, err := DecodeImageFile(path)
	if err != nil {
		return err
	}
	img.Resize(r.width, r.height)
	return saveImage(filepath.Join(r.output, img.Filename()), img)
}

func (r *Resizer) ResizeImage(reader io.Reader, filename string) error {
	img, err := DecodeImage(reader, filename)
	if err != nil {
		return err
	}
	img.Resize(r.width, r.height)
	return saveImage(filepath.Join(r.output, img.Filename()), img)
}

func saveImage(path string, img Image) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer out.Close()

	return img.Encode(out)
}

func scanDir(path string) (<-chan Image, error) {
	imagesCh := make(chan Image)

	// TODO: recursive search by filepath.Walk
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("reading dir %q: %w", path, err)
	}

	go func() {
		defer close(imagesCh)

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			imagePath := filepath.Join(path, file.Name())
			image, err := DecodeImageFile(imagePath)
			if err != nil {
				fmt.Printf("skipping file: %q: %v\n", imagePath, err)
				continue
			}
			imagesCh <- image
		}
	}()
	return imagesCh, nil
}
