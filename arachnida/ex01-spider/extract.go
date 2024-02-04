package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"spider/logger"
	"strings"

	"golang.org/x/net/html"
)

type StringSet map[string]struct{}
var extensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
var visitedUrls = StringSet{}

func isValidImageExtension(url string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

func GetImageName(urlString string) string {
	imageUrl, err := url.Parse(urlString)
    if err != nil {
        logger.LogError(err)
    }

    segments := strings.Split(imageUrl.Path, "/")
    return segments[len(segments)-1]
}

func ExtractImage(url string, path string) bool {
	resp, err := http.Get(url)
    if err != nil {
		logger.LogStr("[error] " + url + ":" + err.Error(), -1)
		return false
    }
	if resp.StatusCode != 200 {
		logger.LogStr("[error] " + url + ":" + resp.Status, -1)
		return false
	}
    defer resp.Body.Close()
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.LogError(err)
	}
	
	imgName := GetImageName(url)
	imgPath := filepath.Join(path, imgName)
	
	if _, err := os.Stat(imgPath); err == nil {
		ext := filepath.Ext(imgName)
		base := strings.TrimSuffix(imgName, ext)
		for i := 1; ; i++ {
			newImgName := fmt.Sprintf("%s_%02d%s", base, i, ext)
			newPath := filepath.Join(path, newImgName)
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				imgPath = newPath
				break
			}
		}
	}
	
	dest, err := os.Create(imgPath)
	if err != nil {
		logger.LogError(err)
    }
	n, err := dest.ReadFrom(resp.Body)
    if err != nil {
		logger.LogError(err)
	} else if n < 0 {
		logger.LogErrorStr("Read failed")
	}
	return true
}

func GetHostname(baseUrl string) string {
	parsedUrl, err := url.Parse(baseUrl)
    if err != nil {
        logger.LogError(err)
    }

    return parsedUrl.Scheme + "://" + parsedUrl.Hostname()
}

func GetValidImageUrl(baseUrl string, imgUrl string) string {
	if strings.HasPrefix(imgUrl, "//") {
		fullUrl := "https:" + imgUrl
		parts := strings.Split(strings.TrimPrefix(imgUrl, "//"), "/")
		if len(parts) > 0 {
			subdomain := len(parts[:len(parts)-1])
			if subdomain > 3 && subdomain < 2 {
				logger.LogStr("[warn] " + fullUrl + ": may be invalid", -1)
			}
		}
		return fullUrl 
	}
	if !strings.HasPrefix(imgUrl, "http")  {
		if strings.HasPrefix(imgUrl, "/") {
			return GetHostname(baseUrl) + imgUrl
		}
		return GetHostname(baseUrl) + "/" + imgUrl
	}
	return imgUrl
}

func ExtractImages(baseUrl string, level int, path string) {
    if level <= 0 {
		return
	}
	resp, err := http.Get(baseUrl)
    if err != nil {
        logger.LogStr("[error] " + baseUrl + ":" + err.Error(), -1)
		return
    }
    defer resp.Body.Close()
    doc, err := html.Parse(resp.Body)
    if err != nil {
        logger.LogError(err)
    }
	
	var foundUrls []string
	imgCount := 0
	var extractImagesFromNodes func(*html.Node)
    extractImagesFromNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" && isValidImageExtension(attr.Val) {
					imgUrl := GetValidImageUrl(baseUrl, attr.Val)
					if ExtractImage(imgUrl, path) {
						imgCount++
					}
                }
            }
        } else if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
					foundUrls = append(foundUrls, attr.Val)
                }
            }
        }
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractImagesFromNodes(c)
		}
		for _, url := range foundUrls {
			if visitedUrls.Contains(url) {
				continue
			}
			visitedUrls.Add(url)
			ExtractImages(url, level - 1, path)
		}
    }
    extractImagesFromNodes(doc)
	logger.LogStr(baseUrl, imgCount)
}
