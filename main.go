package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

var errorQueueEmpty = errors.New("queue is empty")

type queue[T any] struct {
	Urls []T
}

func (q *queue[T]) enqueue(toQueue T) {
	q.Urls = append(q.Urls, toQueue)
}

func (q *queue[T]) dequeue() (T, error) {
	var zero T
	if len(q.Urls) == 0 {
		return zero, errorQueueEmpty
	}
	fmt.Println(len(q.Urls))
	res := q.Urls[0]
	q.Urls[0] = zero
	q.Urls = q.Urls[1:]
	return res, nil
}

func main() {
	baseUrl := "https://web-scraping.dev"
	htmlUrls := map[string][]string{}
	queue := queue[string]{}

	for {
		queue.enqueue(baseUrl)
		url, err := queue.dequeue()
		if err != nil {
			if errors.Is(err, errorQueueEmpty) {
				fmt.Println(err)
				break 
			} else {
				fmt.Println(err)
				return
			}
		}
		urls, err := downloadLink(url, queue)
		if err != nil {
			fmt.Println(err)
			return
		}
		htmlUrls[url] = urls
	}

	// for _, url := range urls {
	// 	link, err := downloadLink(url)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	htmlUrls[url] = link
	// }
	// fmt.Println(htmlUrls)
}

func downloadLink(url string, queue queue[string]) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	token := html.NewTokenizer(resp.Body)
	result := []string{}
	urls := map[string]int{}
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
						toAdd := urlNormalization(url, attribute.Val)
						if toAdd==url{
							continue
						}
						if _, ok := urls[toAdd]; !ok {
							result = append(result, toAdd)
							queue.enqueue(toAdd)
						}
						urls[toAdd] = 1
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

func urlNormalization(baseUrl string, url string) string {
	if strings.Contains(url, "https") {
		return url
	}
	if url == "#" {
		return baseUrl
	} else {
		return baseUrl + url
	}
}
