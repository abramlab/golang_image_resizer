package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/abramlab/resizer/resizer"
)

// TODO: add progress line
// TODO: add tests
// TODO: delete resize pkg
// TODO: add logger
// TODO: add docker
// TODO: add deb pkg
// TODO: rename project (image resizer)
// TODO: update readme
// TODO: add more formats

var (
	input  = flag.String("input", "images", "Path to folder where images you need to resize.")
	output = flag.String("output", "resized_images", "Path to output folder with resized images.")

	width  = flag.Uint("width", 1024, "Width of resized images in px.")
	height = flag.Uint("height", 0, "Height of resized images in px.")

	resFolder = flag.Bool("resolution-folder", false,
		"All resized images will be saved in separate folder. "+
			"For example if width = 1024 and height = 0, resized images will be saved in 'output/1024x0' folder.")
	workersNum = flag.Int("workers_num", runtime.NumCPU(),
		"Number of resized workers. Default is number of logical CPUs.")
)

func main() {
	flag.Parse()

	if *workersNum <= 0 {
		log.Fatalln("number of workers should be > 0")
	}
	if *resFolder {
		*output = filepath.Join(*output, fmt.Sprintf("%dx%d", *width, *height))
	}

	r, err := resizer.NewResizer(
		*input, *output,
		resizer.WithResolution(*width, *height),
		resizer.WithWorkersNum(*workersNum),
	)
	if err != nil {
		log.Fatalf("create new resizer failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		cancel()
		fmt.Println("Stopping signal caught, shutting down...")
	}()

	result, err := r.Run(ctx)
	if err != nil {
		log.Fatalf("run resizer failed: %v", err)
	}

	fmt.Printf("Resized images: %v\n", result.ResizedImages)
}
