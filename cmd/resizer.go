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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	result, err := r.Run(ctx)
	if err != nil {
		log.Fatalf("run resizer failed: %v", err)
	}

	fmt.Printf("RESIZED IMAGES: %v\n", result.ResizedImages)
}
