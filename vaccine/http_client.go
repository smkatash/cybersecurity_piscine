package main

import (
	"errors"
	"fmt"
	"net/http"
	"log"
	"io"
	"net/url"
)


type HttpClient struct {
	baseURL string
	sessionCookie  []*http.Cookie
	verbose bool
}

func (c *HttpClient) get(endPoint string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.baseURL + endPoint, nil)
	if c.verbose == true {
		fmt.Println(req)
	}
    if err != nil {
        log.Fatal(err)
    }
	for _, cookie := range c.sessionCookie {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
	return resp
}

func (c *HttpClient) post(data url.Values) *http.Response {
	resp, err := http.PostForm(c.baseURL, data)
    if err != nil {
		log.Fatal(err)
    }
	if c.verbose == true {
		fmt.Println("POST " + c.baseURL)
		fmt.Println(data)
	}

    return resp
}

func (c *HttpClient) Request(method string, query url.Values) *http.Response {
	if method == "GET" {
		return c.get("?" + query.Encode())
	} else if method == "POST" {
		return c.post(query)
	}
	return nil
}

func (c *HttpClient) Response(response *http.Response) (string, error) {
	if response.StatusCode == 200 {
		rbody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}
		defer response.Body.Close()
		return string(rbody), nil
	}
	return "", errors.New(response.Status)
}


