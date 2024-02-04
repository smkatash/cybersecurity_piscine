package main

import (
	"fmt"
	"os"
	"spider/logger"
)

type Options struct {
	recursive bool
	level int
	path string
	url	string
}

func main() {
	if len(os.Args) < 2 {
		logger.LogUsage()
	}
	fmt.Println(len(os.Args))
	opts := Options {false, 5, "./data", ""}
	ParseInput(os.Args[1:], &opts)
	ExtractImages(opts.url, opts.level, opts.path)
}