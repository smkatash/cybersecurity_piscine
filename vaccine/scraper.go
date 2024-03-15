package main

import (
	"golang.org/x/net/html"
	"strings"
	"fmt"
	"net/url"
)

type Scraper struct {
	client *HttpClient
}

type Form struct {
	method string
	queryValues []string
	queryURL string
	colNum int
}

func (s *Scraper) GetPageValues(inputMethod string) []Form {
	var forms []Form
	response := s.client.Request("GET", url.Values{})
	body, _ := s.client.Response(response)
	formNodes, _ := s.extractHTMLForm(body) 
	fmt.Println("Vulnerable forms found: ")
	for _, form := range formNodes {
		method := s.getFormMethod(form)
		inputNames := s.getFormInputNames(form)
		fmt.Printf("Method: %s\n", method)
		fmt.Printf("Input Names: %v\n", inputNames)
		fmt.Println()
		if strings.ToUpper(inputMethod) != method {
			fmt.Println("Input method is not supported! Proceeding ...")
		}
		forms = append(forms, Form{
			inputMethod,
			inputNames,
			"", 0,
		})
	}
	return forms
}


func (s *Scraper) extractHTMLForm(htmlContent string) ([]*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var formNodes []*html.Node
	var findForms func(*html.Node)
	findForms = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			formNodes = append(formNodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findForms(c)
		}
	}
	findForms(doc)

	return formNodes, nil
}

func (s *Scraper) getFormMethod(formNode *html.Node) string {
	for _, attr := range formNode.Attr {
		if attr.Key == "method" {
			return attr.Val
		}
	}
	return ""
}

func (s *Scraper) getFormInputNames(formNode *html.Node) []string {
	var inputNames []string
	var findInputs func(*html.Node)
	findInputs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			isHidden := false
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "hidden" {
					isHidden = true
					break
				}
			}
			if !isHidden {
				for _, attr := range n.Attr {
					if attr.Key == "name" {
						inputNames = append(inputNames, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findInputs(c)
		}
	}
	findInputs(formNode)

	return inputNames
}
