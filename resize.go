package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
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
	inPath  = flag.String("in_path", "images", "Path to folder where images you need to resize.")
	outPath = flag.String("out_path", "resized", "Path to folder with resized images.")

	width  = flag.Uint("width", 240, "Width of resized images in px.")
	height = flag.Uint("height", 240, "Height of resized images in px.")

	postfix = flag.Bool("postfix", false, "Postfix of width and height in resized image. Example: img_name_300x300.")
	gNum    = flag.Int("gNum", runtime.NumCPU(), "Number of resized workers. Default is number of logical CPUs.")
)

type Resizer struct {
	width, height   uint
	inPath, outPath string
	postfix         bool
	ch              chan Img
	wg              *sync.WaitGroup
	count           int
}

func main() {
	flag.Parse()

	//Error if number of goroutines is <= 0
	if *gNum <= 0 {
		log.Fatalln("number of goroutines should be > 0")
	}

	//Create new resizer with all params
	r := NewResizer(*width, *height, *inPath, *outPath, *postfix)

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
	for i := 0; i < *gNum; i++ {
		r.wg.Add(1)
		go r.resizeWorker()
	}

	r.wg.Wait()
	fmt.Printf("Success resized images: %v\n", r.count)
}

func NewResizer(w, h uint, in, out string, postfix bool) *Resizer {
	return &Resizer{
		width:   w,
		height:  h,
		inPath:  in,
		outPath: out,
		postfix: postfix,
		ch:      make(chan Img),
		wg:      &sync.WaitGroup{},
		count:   0,
	}
}

func (r *Resizer) scanDir() error {
	defer close(r.ch)

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
			r.ch <- Img{img, name, format}
		}
	}
	return nil
}

func (r *Resizer) resizeWorker() {
	defer r.wg.Done()
	for img := range r.ch {
		img.resize(r.width, r.height)
		resizedPath := r.fullOutPath(img)
		if err := img.saveTo(resizedPath); err != nil {
			fmt.Printf("saving resized file error: %v\n", err)
			continue
		}
		r.count++
	}
}

func (r *Resizer) fullOutPath(img Img) string {
	var path string
	if r.postfix {
		w, h := img.dimensions()
		path = fmt.Sprintf("%s/%s_%vx%v%s", r.outPath, img.name, w, h, img.ext())
	} else {
		path = fmt.Sprintf("%s/%s%s", r.outPath, img.name, img.ext())
	}
	return path
}

type Img struct {
	image  image.Image
	name   string
	format string
}

func (i *Img) resize(width, height uint) {
	i.image = resize.Resize(width, height, i.image, resize.Lanczos3)
}

func (i *Img) dimensions() (width int, height int) {
	b := i.image.Bounds()
	return b.Max.X, b.Max.Y
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

func (i Img) ext() string {
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

func decodeImage(path string) (image.Image, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("cant open file %s: %s", path, err)
	}
	defer file.Close()

	bFilename := filepath.Base(file.Name())
	fileName := strings.TrimSuffix(bFilename, filepath.Ext(bFilename))

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", "", fmt.Errorf("cant decode file %s: %s", path, err)
	}

	return img, fileName, format, nil
}
