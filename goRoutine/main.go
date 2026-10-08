package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

type crawlResult struct {
	baseURL string
	urls    []string
	err     error
}

func workers(receiveURL <-chan string, sendURL chan<- crawlResult) {

	url := <-receiveURL
	links, err := downloadLink(url)
	if err != nil {
		sendURL <- crawlResult{url, nil, err}
	}

	sendURL <- crawlResult{url, links, nil}
}

func main() {
	baseURL := "https://web-scraping.dev/"

	receiveURL := make(chan string, 50)
	sendURL := make(chan crawlResult)
	receiveURL <- baseURL
	allURLs := map[string][]string{}

	allLinks := map[string]int{}
	allLinks[baseURL] = 1
	var wg sync.WaitGroup

	maxWorkers := 5
	sem := make(chan struct{}, maxWorkers)
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			workers(receiveURL, sendURL)
			<-sem
		}()
	}

	for urls := range sendURL {
		if urls.err!=nil{
			log.Print(urls.err)
			continue
		}
		toAdd := []string{}
		for _, url := range urls.urls {
			if _, ok := allLinks[url]; !ok {
				allLinks[url] = 1
				toAdd = append(toAdd, url)
				receiveURL<-url
			}
		}
		allURLs[urls.baseURL] = toAdd
	}

	close(receiveURL)
	wg.Wait()
	close(sendURL)
	fmt.Println(allURLs)
}

func downloadLink(baseURL string) ([]string, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	token := html.NewTokenizer(resp.Body)
	result := []string{}

	for {
		tokenType := token.Next()
		switch tokenType {
		case html.ErrorToken:
			if token.Err() == io.EOF {
				fmt.Println("reached end of file")
				return result, nil
			}
			return nil, fmt.Errorf("error: %v", token.Token().Data)
		case html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {

				for _, attribute := range t.Attr {
					if attribute.Key == "href" {
						relativeURL, err := url.Parse(attribute.Val)

						if err != nil {
							continue
						}
						base, err := url.Parse(baseURL)
						if err != nil {
							continue
						}
						toAdd := base.ResolveReference(relativeURL)
						result = append(result, toAdd.String())
					}
				}
			}
		}
	}
}
