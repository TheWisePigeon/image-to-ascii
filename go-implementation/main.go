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
	grayscaleImage := image.NewGray(imageToConvert.Bounds())
	draw.Draw(grayscaleImage, grayscaleImage.Bounds(), imageToConvert, imageToConvert.Bounds().Min, draw.Src)
}
