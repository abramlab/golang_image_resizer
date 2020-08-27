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

## Examples

1.`go run resize.go -width=500 -height=500 -postfix`

| Before | After |
| --- | --- |
| nature1.jpg 800x825 | nature1_500x500.jpg 500x500 |
|![Nature1 image 800x825](test_images/in/nature1.jpg?raw=true)|![Nature1 image 500x500](test_images/out/nature1_500x500.jpg?raw=true)|

2.`go run resize.go -width=600 -height=0`

| Before | After |
| --- | --- |
| nature2.jpg 1366x550 | nature2.jpg 600x242 |
|![Nature2 image 1366x550](test_images/in/nature2.jpg?raw=true)|![Nature2 image 600x242](test_images/out/nature2.jpg?raw=true)|

## Supported flags
    go run resize.go -help
    
The following flags are supported:

| Flag | Default | Description |
| --- | --- | --- |
| `width` | 240 | Width of resized images in px |
| `height` | 240 | Height of resized images in px |
| `in_path` | images | Path to folder where images you need to resize |
| `out_path` | resized | Path to folder with resized images |
| `postfix` | false | Postfix of width and height in resized image. Example: img_name_300x300 |
| `gNum` | Number of logical CPUs | Number of resized workers |