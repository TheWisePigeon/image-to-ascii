package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

var (
	intensityMap = map[int]string{
		0:   " ", // Lowest intensity, e.g., background
		50:  ".", // Low intensity
		100: "-", // Moderate intensity
		150: "+", // Higher intensity
		200: "*", // Even higher intensity
		255: "@", // Highest intensity, e.g., foreground
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Wrong usage. Missing file path")
		os.Exit(1)
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Invalid path:", filePath, "is not a file")
			os.Exit(1)
		}
		panic(err)
	}
	defer file.Close()
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}
	fileType := http.DetectContentType(buffer)
	if fileType != "image/png" && fileType != "image/jpeg" {
		fmt.Println("Format not supported")
		os.Exit(1)
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	var imageToConvert image.Image
	if fileType == "image/png" {
		imageToConvert, err = png.Decode(file)
		if err != nil {
			panic(err)
		}
	}
	if fileType == "image/jpeg" {
		imageToConvert, err = jpeg.Decode(file)
		if err != nil {
			panic(err)
		}
	}
	resized := resize.Resize(10, 0, imageToConvert, resize.NearestNeighbor)
	grayscaleImage := image.NewGray(resized.Bounds())
	draw.Draw(grayscaleImage, grayscaleImage.Bounds(), imageToConvert, imageToConvert.Bounds().Min, draw.Src)
	result := ""
	for y := 0; y < grayscaleImage.Bounds().Max.Y; y++ {
		for x := 0; x < grayscaleImage.Bounds().Max.X; x++ {
			r, _, _, _ := grayscaleImage.At(x, y).RGBA()
			intensity := int(r >> 8)
			if intensity == 0 {
				result += " "
				continue
			}
			if intensity >= 0 && intensity < 50 {
				result += "."
				continue
			}
			if intensity >= 50 && intensity < 100 {
				result += "-"
				continue
			}
			if intensity >= 100 && intensity < 150 {
				result += "+"
				continue
			}
			if intensity >= 150 && intensity < 200 {
				result += "*"
				continue
			}
			if intensity >= 200 {
				result += "@"
				continue
			}
		}
		result += "\n"
	}
	fmt.Println(result)
}
