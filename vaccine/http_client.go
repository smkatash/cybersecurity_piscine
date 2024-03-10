package main

import (
	"errors"
	"fmt"
	"net/http"
	"log"
	"io"
)


type HttpClient struct {
	baseURL string
	sessionCookie  []*http.Cookie
}

func (c *HttpClient) Get(endPoint string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.baseURL + endPoint, nil)
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

// func (c *HttpClient) Post(endPoint string, payload string) *http.Response {
// 	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	client := &http.Client{}
// 	fmt.Println("Sending to ", endPoint)
// 	req, err := http.NewRequest("POST", endPoint, &payload)
//     if err != nil {
//         log.Fatal(err)
//     }
//     req.AddCookie(&sessionCookie)
// 	resp, err := client.Do(req)
//     if err != nil {
//         log.Fatal(err)
//     }
// 	return resp
// }

func (c *HttpClient) Response(response *http.Response) (string, error) {
	if response.StatusCode == 200 {
		rbody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return string(rbody), nil
	}
	return "", errors.New(response.Status)
}


