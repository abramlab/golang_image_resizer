## Golang image resizer
This script is used for resizing images( jpg | png | gif ).

## Installation
First, if don`t have Go, install it - [instruction](https://golang.org/doc/install#install)
	
	git clone https://github.com/abram213/golang_image_resizer.git
	cd golang_image_resizer/

## Running	
1. Drop images that need to resize into `images` folder
2. Run script `go run resize.go -width=500 -height=500`
3. By default, all resized images saves in `resized` folder

## Supported flags
    go run resize.go -help
    
The following flags are supported:

| Flag | Default | Description |
| --- | --- | --- |
| `width` | 240 | Width of resized images in px |
| `height` | 240 | Height of resized images in px |
| `in_path` | images | Path to folder where images you need to resize |
| `out_path` | resized | Path to folder with resized images |
| `gNum` | Number of logical CPUs | Number of resized workers |