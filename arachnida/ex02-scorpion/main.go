package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var extensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

// TODO bmp format!
func IsValidImage(img string) bool {
	for _, ext := range extensions {
	  if strings.HasSuffix(strings.ToLower(img), ext) {
		return true
	  }
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("scorpion: Usage ./scorpion FILE1 [FILE2 ...]")
	}

	for i, arg := range os.Args[1:] {
		fmt.Printf("%d.Image: %s\n", i + 1, arg)
		if !IsValidImage(arg) {
			log.Fatal("scorpion: Invalid image file")
		}
		Decode(arg)
	}
}