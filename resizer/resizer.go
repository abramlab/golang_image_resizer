package resizer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
	imagesCh, err := scanInputDir(r.input)
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
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %s: %s", path, err)
	}
	defer file.Close()

	img, err := DecodeImage(file, filepath.Base(file.Name()))
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
	if err := createDirForFile(path, 0o755); err != nil {
		return fmt.Errorf("creating directory for file %q: %w", path, err)
	}
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %q: %w", path, err)
	}
	defer out.Close()

	return img.Encode(out)
}

func scanInputDir(dir string) (<-chan Image, error) {
	imagesCh := make(chan Image)

	go func() {
		defer close(imagesCh)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("reading file: %v\n", err)
				return nil
			}
			if info.IsDir() {
				return nil
			}
			image, err := decodeImageFile(path, dir+string(filepath.Separator))
			if err != nil {
				fmt.Printf("skipping file: %q: %v\n", path, err)
				return nil
			}
			imagesCh <- image
			return nil
		})
		if err != nil {
			fmt.Printf("walking directory %q: %v", dir, err)
		}
	}()
	return imagesCh, nil
}

func createDirForFile(path string, dirPerm os.FileMode) error {
	if base := filepath.Base(path); base == "." || base == "/" {
		return errors.New("path doesn't point to file")
	}
	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}
	return os.MkdirAll(dir, dirPerm)
}
