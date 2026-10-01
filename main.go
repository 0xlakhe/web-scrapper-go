package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	// "strings"
)

func main() {
	urls := []string{"https://web-scraping.dev/"}
	htmlUrls := map[string]map[int]string{}

	for _, url := range urls {
		link, err := downloadLink(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		htmlUrls[url] = link
	}
	for key, value := range htmlUrls {
		fmt.Printf("%q\n", key)
		for k, v := range value {
			fmt.Println(k, v)
		}
	}
}

func downloadLink(url string) (map[int]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	token := html.NewTokenizer(resp.Body)
	result := map[int]string{}
	idx := 1
outer:
	for {
		tokenType := token.Next()
		switch tokenType {
		case html.ErrorToken:
			return nil, fmt.Errorf("error: %v", token.Token().Data)
		case html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {

				for _, attribute := range t.Attr {
					if attribute.Key == "href" {
						result[idx] = attribute.Val
						idx += 1
					}
				}
			}
		case html.EndTagToken:
			if token.Token().Data == "html" {
				break outer
			}
		}
	}
	return result, nil
}
