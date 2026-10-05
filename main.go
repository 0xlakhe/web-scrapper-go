package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type queue[T any] struct{
	Urls []T
}

func (q *queue[T]) enqueue(toQueue T) {
	q.Urls = append(q.Urls, toQueue)
}

func (q *queue[T]) dequeue() (T,error){
	var zero T
	if len(q.Urls)==0{
		return zero,fmt.Errorf("queue is empty")
	}
	fmt.Println(len(q.Urls))
	res:=q.Urls[0]
	q.Urls[0]=zero
	q.Urls=q.Urls[1:]
	return res,nil
}


func main() {
	urls:=[]string{"https://web-scraping.dev"}
	htmlUrls:=map[string][]string{}

	for _, url := range urls {
		link, err := downloadLink(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		htmlUrls[url] = link
	}
	fmt.Println(htmlUrls)
}

func downloadLink(url string)([]string,error){
	resp,err:=http.Get(url)
	if err!=nil{
		return nil,err
	}
	token:=html.NewTokenizer(resp.Body) 
	result:=[]string{}
	urls:=map[string]int{}
	outer:
	for{
		tokenType:=token.Next()
		switch tokenType{
		case html.ErrorToken:
			return nil, fmt.Errorf("error: %v", token.Token().Data)
		case html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {

				for _,attribute:=range t.Attr{
					if attribute.Key=="href"{
						toAdd:=urlNormalization(url,attribute.Val)

						if _,ok:=urls[toAdd]; !ok{
							result=append(result,toAdd)
						}
						urls[toAdd]=1
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

func urlNormalization(baseUrl string,url string) string{
		if strings.Contains(url,"https"){
			return url
		}
		if url=="#"{
			return baseUrl+"/"+url
		}else{
			return baseUrl+url
		}
}