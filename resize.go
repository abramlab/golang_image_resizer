package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"gocv.io/x/gocv"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	inPath = flag.String("img_path", "images", "")
	outPath = flag.String("out_path", "resized", "")

	width = flag.Uint("width", 240, "Width of resized image")
	height = flag.Uint("height", 240, "Height of resized image")

	gCount = flag.Int("g", runtime.NumCPU(), "Number of working goroutines")
)

type Resizer struct {
	width, height uint
	inPath, outPath string
	ch chan Img
	wg *sync.WaitGroup
	count int
}

func main() {
	flag.Parse()

	//Error if number of goroutines is <= 0
	if *gCount <= 0 {
		log.Fatalln("number of goroutines should be > 0")
	}

	//Create new resizer with all params
	r := NewResizer(*width, *height, *inPath, *outPath)

	//If out dir doesn't exist, create it
	resizedPath := *outPath
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		os.MkdirAll(resizedPath, os.ModePerm)
	}

	//Start scanning dir with images in separate goroutine
	go func() {
		if err := r.scanDir(); err != nil {
			log.Fatalf("scanning dir error: %v", err)
		}
	}()

	//Start resize workers
	for i := 0; i < *gCount; i++ {
		r.wg.Add(1)
		go r.resizeWorker()
	}

	r.wg.Wait()
	fmt.Printf("Success resized images: %v\n", r.count)
}

func NewResizer(w, h uint, in, out string) *Resizer {
	return &Resizer{
		width: 		w,
		height:     h,
		inPath: 	in,
		outPath:	out,
		ch: 		make(chan Img),
		wg: 		&sync.WaitGroup{},
		count: 		0,
	}
}

func (r *Resizer) scanDir() error {
	dirItems, err := ioutil.ReadDir(r.inPath)
	if err != nil {
		return errors.Wrap(err, "read dir error")
	}

	for _, item := range dirItems {
		if !item.IsDir() {
			imgPath := fmt.Sprintf("%s/%s", r.inPath, item.Name())
			img, name, format, err := decodeImage(imgPath)
			if err != nil {
				fmt.Printf("skipping file: %v\n", err)
				continue
			}
			r.ch <- Img{img, name,format}
		}
	}
	close(r.ch)
	return nil
}

func (r *Resizer) resizeWorker() {
	for img := range r.ch {
		img.resize(r.width, r.height)
		resizedPath := fmt.Sprintf("%s/%s_%vx%v%s", r.outPath, img.name, r.width, r.height, img.ext())
		if err := img.saveTo(resizedPath); err != nil {
			fmt.Printf("saving resized file error: %v\n", err)
			continue
		}
		r.count++
	}
	r.wg.Done()
}

type Img struct {
	image 	image.Image
	name 	string
	format 	string
}

func (i *Img) resize(width, height uint) {
	i.image = resize.Resize(width, height, i.image, resize.Lanczos3)
}

func (i Img) saveTo(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cant create file %s: %s", path, err)
	}
	defer out.Close()

	switch i.format {
	case "jpeg":
		return jpeg.Encode(out, i.image, nil)
	case "png":
		return png.Encode(out, i.image)
	case "gif":
		return gif.Encode(out, i.image, nil)
	default:
		return fmt.Errorf("unsupported format: %s", i.format)
	}
}

func (i Img) ext() gocv.FileExt {
	var ext gocv.FileExt
	switch i.format {
	case "jpeg":
		ext = gocv.JPEGFileExt
	case "png":
		ext = gocv.PNGFileExt
	case "gif":
		ext = gocv.GIFFileExt
	}
	return ext
}

func decodeImage(path string) (image.Image, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("cant open file %s: %s", path, err)
	}
	bFilename := filepath.Base(file.Name())
	fileName := strings.TrimSuffix(bFilename, filepath.Ext(bFilename))

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", "", fmt.Errorf("cant decode file %s: %s", path, err)
	}
	file.Close()
	return img, fileName, format, nil
}