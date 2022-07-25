package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/abramlab/resizer/resizer"
)

// TODO: add progress line
// TODO: add tests
// TODO: in out path can be also files
// TODO: delete resize pkg
// TODO: add logger
// TODO: add context
// TODO: add makefile
// TODO: add docker
// TODO: add deb pkg
// TODO: rename project (image resizer)
// TODO: update readme
// TODO: add more formats

var (
	inPath  = flag.String("in_path", "images", "Path to folder where images you need to resize.")
	outPath = flag.String("out_path", "resized_images", "Path to folder with resized images.")

	width  = flag.Uint("width", 240, "Width of resized images in px.")
	height = flag.Uint("height", 240, "Height of resized images in px.")

	// TODO: for future use
	/*postfix = flag.Bool("postfix", false,
	"Postfix of width and height in resized image. Example: img_name_300x300.")*/
	workersNum = flag.Int("workers_num", runtime.NumCPU(),
		"Number of resized workers. Default is number of logical CPUs.")
)

func main() {
	flag.Parse()

	if *workersNum <= 0 {
		log.Fatalln("number of workers should be > 0")
	}

	r, err := resizer.NewResizer(*outPath, resizer.WithResolution(*width, *height))
	if err != nil {
		log.Fatalf("create new resizer failed: %v", err)
	}

	imagesCh, err := scanDir(*inPath)
	if err != nil {
		log.Fatalf("scanning dir failed: %v", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < *workersNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.RunResizeWorker(imagesCh)
		}()
	}

	wg.Wait()
	// fmt.Printf("Success resized images: %v\n", r.resizedImages)
}

func scanDir(path string) (<-chan resizer.Image, error) {
	imagesCh := make(chan resizer.Image)

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
			image, err := resizer.DecodeImageFile(filepath.Join(path, file.Name()))
			if err != nil {
				fmt.Printf("skipping file: %v\n", err)
				continue
			}
			imagesCh <- image
		}
	}()
	return imagesCh, nil
}
