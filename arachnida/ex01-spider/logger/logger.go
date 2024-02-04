package logger

import (
	"log"
	"fmt"
)

func LogUsage() {
	log.Fatal("usage: ./spider [-rlp] URL")
}


func LogError(err error) {
	log.Fatalf("spider: %s", err)
}

func LogErrorStr(err string) {
	log.Fatalf("spider: %s", err)
}

func LogStr(str string, d int) {
	if d >= 0 {
		fmt.Printf("spider: %s %d\n", str, d)
	} else {
		fmt.Printf("spider: %s\n", str)
	}
}