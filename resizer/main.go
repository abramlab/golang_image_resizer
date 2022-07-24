package resizer

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

// TODO: add enums
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

	postfix = flag.Bool("postfix", false,
		"Postfix of width and height in resized image. Example: img_name_300x300.")
	workersNum = flag.Int("workers_num", runtime.NumCPU(),
		"Number of resized workers. Default is number of logical CPUs.")
)

const publicFolderMode = 0o755

func main() {
	flag.Parse()

	if *workersNum <= 0 {
		log.Fatalln("number of workers should be > 0")
	}

	r := newResizer(*width, *height, *inPath, *outPath, *postfix)

	out := *outPath
	if err := os.MkdirAll(out, publicFolderMode); err != nil {
		log.Fatalf("can't create %q dir: %v", out, err)
	}

	go func() {
		if err := r.scanDir(); err != nil {
			log.Fatalf("scanning dir error: %v", err)
		}
	}()

	for i := 0; i < *workersNum; i++ {
		r.wg.Add(1)
		go r.resizeWorker()
	}

	r.wg.Wait()
	fmt.Printf("Success resized images: %v\n", r.resizedImages)
}
