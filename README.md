## Image resizer

Easily and quickly resize local images( `jpg` | `png` | `gif` ).

## Running

Drag the images you want to resize to the `images`(default) folder.
All resized images are saved in the `resized-images`(default) folder.

#### Binary
```
image-resizer -width=500 -height=500
```

#### Docker image

```
docker run \
    -it \
    --rm \
    -w /app \
    --mount type=bind,source=$(pwd)/images,target=/app/images \
    --mount type=bind,source=$(pwd)/resized-images,target=/app/resized-images \
    abramlab/image-resizer:0.1.0
```

## Build

First, if you don`t have Go, [install](https://golang.org/doc/install#install) it.

1. `git clone https://github.com/abramlab/image-resizer.git`
2. `cd image-resizer`
3. `make build` or `go build -o bin/image-resizer .`

## Supported options

```
  --input=                      Path to folder where images you need to resize. (default: ./images)
  --output=                     Path to output folder with resized images. (default: ./resized_images)
  --width=                      Width of resized images in px. (default: 1024)
  --height=                     Height of resized images in px. (default: 0)
  --resolution_folder           All resized images will be saved in separate folder.
		                For example if width = 1024 and height = 0, resized images will be saved 
		                in 'output/1024x0' folder. (default: false)
  --workers_num=                Number of resized workers. (default: number of logical CPUs)
```