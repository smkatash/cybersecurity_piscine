package main

import (
	"net/url"
	"os"
	"spider/logger"
	"strconv"
	"strings"
)

func IsValidUrl(urlString string) bool {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}


func ParseInput(args []string, opts *Options) {
	len := len(args)
	if !IsValidUrl(args[len - 1]) {
		logger.LogErrorStr("Invalid url")
	}
	opts.url = args[len - 1]
	len--
	
	for x := 0; x < len; x++ {
		switch args[x] {
		case "-r":
			opts.recursive = true
		case "-l":
			if x + 1 < len {
				if num, err := strconv.Atoi(args[x + 1]); err == nil {
					opts.level = num
					x++
				} else {
					if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
						logger.LogError(numErr.Err)
					}
				}
			}
		case "-p":
			if x + 1 < len {
				_, err := os.Stat(args[x + 1])
				if err == nil {
					opts.path = strings.TrimRight(args[x + 1], "/")
					x++
				} else if os.IsNotExist(err) {
					logger.LogError(err)
				}
			}
		default:
			logger.LogUsage()
	  }
	}
	
}