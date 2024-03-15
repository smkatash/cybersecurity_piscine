package main

import (
	"net/http"
	"log"
	"net/url"
	"os"
)

type Input struct {
	filename string
	method string
	baseURL string
	verbose bool
	session string
	sessionCookie  []*http.Cookie
}

func Start(input *Input) {
	client := HttpClient{input.baseURL, input.sessionCookie, input.verbose}
	scraper := Scraper{&client}
	forms := scraper.GetPageValues(input.method)
	injector := SqlInjector{&client, forms, input.baseURL, "", "--", "", "", ""}
	if err := injector.Initialize(); err != nil {
        log.Fatal("Error initializing SQL injector: ", err)
    }
	if err := injector.GetDatabaseDump(input.filename); err != nil {
        log.Fatal("Error getting database dump: ", err)
    }
}

func ParseInput(input *Input) {
	args := os.Args
	var flagIndex int
	usageMsg := "usage: ./vaccine [-osvX] URL"
	if len(args) <= 1 {
		log.Fatal(usageMsg)
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-o", "-X", "-s":
			flagIndex = i
			i++
			if i >= len(args) {
				log.Fatal(usageMsg)
			}
			switch args[flagIndex] {
			case "-o":
				input.filename = args[i]
				if _, err := os.Stat(input.filename); os.IsNotExist(err) {
					log.Fatal(err)
				}
			case "-X":
				input.method = args[i]
			case "-s":
				input.session = args[i]
			}
		case "-v":
			input.verbose = true

		default:
			if len(input.baseURL) > 0 {
				log.Fatal(usageMsg)
			} 
			u, err := url.Parse(args[i]); if err != nil {
				log.Fatal(usageMsg)
			}
			if u.Scheme == "" {
				log.Fatal("Unsupported protocol scheme for: ", args[i])
			}
	
			input.baseURL = args[i]
		}
	}
}

func main() {
	input := Input{"db_dump.txt", "GET", "", false, "", nil}
	ParseInput(&input) 
	if len(input.session) > 0 {
		sessionCookie := []*http.Cookie{
			{Name: "PHPSESSID", Value: input.session},
			{Name: "security", Value: "low"},
		}
		input.sessionCookie = sessionCookie
	}
	Start(&input)
}
