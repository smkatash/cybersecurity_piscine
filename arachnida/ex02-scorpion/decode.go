package main

import (
	"fmt"
	"log"
	"os"
	"image"
	_ "image/jpeg"
    _ "image/png"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type Exif struct{}

func (p Exif) Walk(name exif.FieldName, tag *tiff.Tag) error {
    fmt.Printf("%s: %s\n", name, tag)
    return nil
}

func DecodeExif(img *os.File) {
	metaData, err := exif.Decode(img)
	if err != nil {
		log.Println("scorpion: exif cannot be decoded: " + err.Error())
		return
	}

	var data Exif
    metaData.Walk(data)
}

func DecodeImage(img *os.File) {
	imgData, _, err := image.Decode(img)
    if err != nil {
        log.Println("scorpion: image cannot be decoded: " + err.Error())
		return
    }

    bounds := imgData.Bounds()
	fmt.Println("Image Width:", bounds.Max.X)
    fmt.Println("Image Height:", bounds.Max.Y)
}

func Decode(imgFile string) {
	img, err := os.Open(imgFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer img.Close()
	DecodeImage(img)
	img.Seek(0, 0)
	DecodeExif(img)
}