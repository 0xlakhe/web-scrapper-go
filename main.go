package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"net/http"
	"net/url"

	"golang.org/x/net/html"
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
	baseUrl := "https://web-scraping.dev/"
	htmlUrls := map[string][]string{}
	queue := queue[string]{}
	queue.enqueue(baseUrl)
	allUrls := map[string]int{}
	allUrls[baseUrl] = 1
	maxPages:=50
	currentPage:=0
	for  currentPage<maxPages {
		baseURL, err := queue.dequeue()
		if err != nil {
			if errors.Is(err, errorQueueEmpty) {
				fmt.Println(err)
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		urls, err := downloadLink(baseURL)

		allLinks:=[]string{}
		//adding links to map
		for _,link:=range urls{
			
			//adding link to queue
			if _,ok:=allUrls[link];!ok{
				allLinks=append(allLinks, link)
				queue.enqueue(link)
			}
			allUrls[link]=1
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\n\n, %v", htmlUrls)

		//saving urls 
		htmlUrls[baseURL] = allLinks
		time.Sleep(2 * time.Second)
		currentPage+=1
	}

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
						result=append(result,toAdd.String())
					}
				}
			}
		}
	}
}
